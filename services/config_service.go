package services

import (
	"errors"
	"fmt"
	"quiver/database"
	"quiver/models"
	"quiver/utils"
	"time"

	"github.com/vmihailenco/msgpack/v5"
	"gorm.io/gorm"
)

// ConfigService 配置服务
type ConfigService struct{}

// NewConfigService 创建配置服务实例
func NewConfigService() *ConfigService {
	return &ConfigService{}
}

// ConfigResponse 配置响应结构
type ConfigResponse struct {
	Env           string            `json:"env"`
	AppName       string            `json:"appName"`
	ClusterName   string            `json:"clusterName"`
	NamespaceName string            `json:"namespaceName"`
	ReleaseKey    string            `json:"releaseKey"`
	Comment       string            `json:"comment"`
	Total         int64             `json:"total"`
	Page          int               `json:"page"`
	Size          int               `json:"size"`
	Items         map[string]string `json:"items"`
	Changes       *Changes          `json:"changes,omitempty"`
}

// Changes 配置变更信息
type Changes struct {
	Updates []string `json:"updates"`
	Adds    []string `json:"adds"`
	Deletes []string `json:"deletes"`
}

// GetNamespaceConfigs 获取命名空间配置
func (s *ConfigService) GetNamespaceConfigs(env, appID, clusterName, namespaceName, releaseKey string, page, size int) (*ConfigResponse, error) {
	// 先从缓存获取

	// 缓存未命中，从数据库查询
	return s.getConfigsFromDB(env, appID, clusterName, namespaceName, releaseKey, page, size)
}

// getConfigsFromDB 从数据库获取配置
func (s *ConfigService) getConfigsFromDB(env string, appName, clusterName, namespaceName, releaseKey string, page, size int) (*ConfigResponse, error) {
	db := database.GetDB(env)

	// 获取最新发布记录
	var release models.NamespaceRelease
	err := db.Where("namespace_name = ?", namespaceName).
		Order("release_time DESC").
		First(&release).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no release found")
		}
		return nil, err
	}

	// 获取发布的配置项
	var itemReleases []models.ItemRelease
	err = db.Where("release_id = ?", release.ReleaseID).Find(&itemReleases).Error
	if err != nil {
		return nil, err
	}

	// 转换为map
	items := make(map[string]string)
	for _, item := range itemReleases {
		items[item.K] = item.V
	}

	response := &ConfigResponse{
		Env:           env,
		AppName:       appName,
		ClusterName:   clusterName,
		NamespaceName: namespaceName,
		ReleaseKey:    release.ReleaseID,
		Comment:       release.Comment,
		Items:         items,
		Total:         int64(len(items)),
		Page:          page,
		Size:          size,
	}

	// 计算变更
	if releaseKey != "" && releaseKey != release.ReleaseID {
		changes, err := s.calculateChanges(env, namespaceName, releaseKey, release.ReleaseID)
		if err == nil {
			response.Changes = changes
		}
	}

	return response, nil
}

// calculateChanges 计算配置变更
func (s *ConfigService) calculateChanges(env string, namespaceName, oldReleaseKey, newReleaseKey string) (*Changes, error) {
	db := database.GetDB(env)

	// 获取旧版本配置
	var oldItems []models.ItemRelease
	err := db.Where("release_id = ?", oldReleaseKey).Find(&oldItems).Error
	if err != nil {
		return nil, err
	}

	// 获取新版本配置
	var newItems []models.ItemRelease
	err = db.Where("release_id = ?", newReleaseKey).Find(&newItems).Error
	if err != nil {
		return nil, err
	}

	// 转换为map便于比较
	oldMap := make(map[string]string)
	for _, item := range oldItems {
		oldMap[item.K] = item.V
	}

	newMap := make(map[string]string)
	for _, item := range newItems {
		newMap[item.K] = item.V
	}

	changes := &Changes{
		Updates: []string{},
		Adds:    []string{},
		Deletes: []string{},
	}

	// 查找新增和更新
	for key, newValue := range newMap {
		if oldValue, exists := oldMap[key]; exists {
			if oldValue != newValue {
				changes.Updates = append(changes.Updates, key)
			}
		} else {
			changes.Adds = append(changes.Adds, key)
		}
	}

	// 查找删除
	for key := range oldMap {
		if _, exists := newMap[key]; !exists {
			changes.Deletes = append(changes.Deletes, key)
		}
	}

	return changes, nil
}

// ReleaseNamespace 发布命名空间
func (s *ConfigService) ReleaseNamespace(env string, namespaceName, operator, comment string) (error, error) {
	db := database.GetDB(env)

	var release *models.NamespaceRelease

	err := db.Transaction(func(tx *gorm.DB) error {
		// 获取当前命名空间的所有配置项
		var items []models.Item
		if err := tx.Where("namespace_name = ?", namespaceName).Find(&items).Error; err != nil {
			return err
		}

		// 生成发布ID
		releaseID := fmt.Sprintf("%s-%d", namespaceName, time.Now().Unix())

		// 创建发布记录
		release = &models.NamespaceRelease{
			NamespaceName: namespaceName,
			ReleaseID:     releaseID,
			Operator:      operator,
			Comment:       comment,
		}

		// 序列化配置数据
		configData := make(map[string]interface{})
		itemsMap := make(map[string]string)

		for _, item := range items {
			kvID := utils.MurmurHash64(item.K, item.V)
			configData[fmt.Sprintf("%d", kvID)] = map[string]interface{}{
				"k": item.K,
				"v": item.V,
			}
			itemsMap[item.K] = item.V

			// 创建配置项发布记录
			itemRelease := models.ItemRelease{
				ReleaseID:   releaseID,
				K:           item.K,
				V:           item.V,
				NamespaceID: item.NamespaceID,
				KvID:        uint64(kvID),
			}

			if err := tx.Create(&itemRelease).Error; err != nil {
				return err
			}
		}

		// 使用MessagePack编码配置数据
		configBytes, err := msgpack.Marshal(configData)
		if err != nil {
			return err
		}
		release.Config = configBytes

		// 保存发布记录
		if err := tx.Create(release).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err, nil
	}

	// 可选：在这里可以使用 release 做一些后续处理，如缓存等
	_ = release
	return nil, nil
}
