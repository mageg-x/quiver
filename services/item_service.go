package services

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"quiver/database"
	"quiver/logger"
	"quiver/models"
	"quiver/utils"
	"strings"
)

// ItemService 命名空间服务
type ItemService struct{}

// NewItemService 创建命名空间服务实例
func NewItemService() *ItemService {
	return &ItemService{}
}

func BulkMarkDeleted[T models.Item | models.ItemRelease](tx *gorm.DB, conditions map[string]interface{}) error {
	for {
		result := tx.Model(new(T)).Where(conditions).Limit(1000).Update("is_deleted", 1)

		if result.Error != nil {
			logger.GetLogger("quiver").Errorf("batch update %T failed %s", *new(T), result.Error)
			return result.Error
		}

		if result.RowsAffected == 0 {
			logger.GetLogger("quiver").Errorf("batch updated 0 records for %T", *new(T))
			break
		}

		logger.GetLogger("quiver").Infof("batch updated %d records for %T", result.RowsAffected, *new(T))
	}
	return nil
}

// CheckACNKFormat 检查 app, cluster, namespace,itemKey 是否合法
func CheckACNKFormat(c *fiber.Ctx, env, appName, clusterName, namespaceName, itemKey *string) (bool, error) {
	if env == nil {
		return true, nil
	}

	if !utils.ValidateEnv(*env) {
		logger.GetLogger("quiver").Errorf("invalid env %s", *env)
		return false, utils.BadRequest(c, "invalid env")
	}

	if appName == nil {
		return true, nil
	}
	if !utils.ValidateAppName(*appName) {
		logger.GetLogger("quiver").Errorf("invalid app_name %s", *appName)
		return false, utils.BadRequest(c, "invalid app_name")
	}

	if clusterName == nil {
		return true, nil
	}

	if !utils.ValidateClusterName(*clusterName) {
		logger.GetLogger("quiver").Errorf("invalid cluster_name %s", *clusterName)
		return false, utils.BadRequest(c, "invalid cluster_name")
	}

	if namespaceName == nil {
		return true, nil
	}
	if !utils.ValidateNamespaceName(*namespaceName) {
		logger.GetLogger("quiver").Errorf("invalid namespace_name %s", *namespaceName)
		return false, utils.BadRequest(c, "invalid namespace_name")
	}

	if itemKey == nil {
		return true, nil
	}
	if !utils.ValidateItemKey(*itemKey) {
		logger.GetLogger("quiver").Errorf("invalid item_key %s", *itemKey)
		return false, utils.BadRequest(c, "invalid item_key")
	}
	return true, nil
}

// CheckACNKinDB 检查 app, cluster, namespace,itemKey 是否在数据库中存在
func CheckACNKinDB(env, appName, clusterName, namespaceName, itemKey *string) (*models.IDs, error) {
	if env == nil || *env == "" {
		logger.GetLogger("quiver").Errorf("env not found")
		return nil, errors.New("env not found")
	}

	var ids models.IDs
	db := database.GetDB(*env)

	if appName == nil {
		return &ids, nil
	}

	// 检查应用是否存在
	var app models.App
	err := db.Where("app_name = ?", *appName).First(&app).Error
	if err != nil {
		logger.GetLogger("quiver").Errorf("app %s not found", *appName)
		return nil, errors.New("app not found")
	}
	ids.AppID = app.AppID

	if clusterName == nil {
		return &ids, nil
	}

	// 检查集群是否存在
	var cluster models.Cluster
	err = db.Where("app_id = ? AND cluster_name = ?", app.AppID, *clusterName).First(&cluster).Error
	if err != nil {
		logger.GetLogger("quiver").Errorf("cluster %s not found", *clusterName)
		return nil, errors.New("cluster not found")
	}
	ids.ClusterID = cluster.ClusterID

	if namespaceName == nil {
		return &ids, nil
	}

	// 检查命名空间是否存在
	var namespace models.Namespace
	err = db.Where("cluster_id = ? AND namespace_name = ?", cluster.ClusterID, *namespaceName).First(&namespace).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.GetLogger("quiver").Errorf("namespace %s not found in ClusterID %d", *namespaceName, cluster.ClusterID)
			return nil, errors.New("namespace not found")
		}
		return nil, errors.New("namespace not found")
	}
	ids.NamespaceID = namespace.NamespaceID

	if itemKey == nil {
		return &ids, nil
	}

	// 检查配置项是否存在
	var item models.Item
	err = db.Where("namespace_id = ? AND k = ?", namespace.NamespaceID, *itemKey).First(&item).Error
	if err != nil {
		logger.GetLogger("quiver").Errorf("item %s not found", *itemKey)
		return nil, errors.New("item not found")
	}
	ids.ItemID = item.ItemID

	return &ids, nil
}

// SetItem 设置配置项
func (s *ItemService) SetItem(env string, appName, clusterName, namespaceName, key, value string) error {
	// 检查 app, cluster, namespace
	ids, err := CheckACNKinDB(&env, &appName, &clusterName, &namespaceName, nil)
	if err != nil {
		return err
	}

	db := database.GetDB(env)
	// 更新或创建配置项
	var item models.Item
	err = db.Where("namespace_id = ? AND k = ?", ids.NamespaceID, key).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			kvID := utils.MurmurHash64(key, value)
			// 创建新配置项
			item = models.Item{
				AppID:       ids.AppID,
				ClusterID:   ids.ClusterID,
				NamespaceID: ids.NamespaceID,
				K:           key,
				V:           value,
				KVId:        uint64(kvID),
				IsReleased:  0,
			}
			err = db.Create(&item).Error
		} else {
			logger.GetLogger("quiver").Errorf("failed to create item")
			return err
		}
	} else {
		// 更新现有配置项
		item.V = value
		item.IsReleased = 0
		item.KVId = uint64(utils.MurmurHash64(key, value))

		err = db.Model(&item).Updates(map[string]interface{}{
			"v":           value,
			"is_released": 0,
			"kv_id":       item.KVId,
		}).Error
	}

	if err != nil {
		logger.GetLogger("quiver").Errorf("failed to update item")
		return err
	}

	return nil
}

// GetItem 获取单个配置项
func (s *ItemService) GetItem(env, appName, clusterName, namespaceName, key string) (*models.Item, error) {
	// 检查 app, cluster, namespace
	ids, err := CheckACNKinDB(&env, &appName, &clusterName, &namespaceName, &key)
	if err != nil {
		return nil, err
	}

	// 检查配置项是否存在
	db := database.GetDB(env)
	var item models.Item
	err = db.Where("namespace_id = ? AND k = ?", ids.NamespaceID, key).First(&item).Error
	if err != nil {
		logger.GetLogger("quiver").Errorf("item %s not found", key)
		return nil, errors.New("item not found")
	}

	return &item, nil
}

func (s *ItemService) ListItem(env, appName, clusterName, namespaceName, search string, page, size int) ([]models.Item, int64, error) {
	// 检查 app, cluster, namespace
	ids, err := CheckACNKinDB(&env, &appName, &clusterName, &namespaceName, nil)
	if err != nil {
		return nil, 0, err
	}
	db := database.GetDB(env)

	var items []models.Item
	var total int64

	// 构建查询条件
	query := db.Where("namespace_id = ? AND is_deleted = 0", ids.NamespaceID)

	// 添加搜索条件（如果 search 参数不为空）
	if search != "" {
		// 使用 LOWER() 实现不区分大小写的模糊匹配
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(`k`) LIKE ?", searchPattern)
	}

	// 获取总数（用于分页）
	if err := query.Model(&models.Item{}).Count(&total).Error; err != nil {
		logger.GetLogger("quiver").Errorf("failed to count items num: %s for namespace :%s", err.Error(), namespaceName)
		return nil, 0, fmt.Errorf("failed to count items: %w", err)
	}

	offset := (page - 1) * size
	query = query.Order("update_time DESC, id DESC").Offset(offset).Limit(size)

	// 执行分页查询
	if err := query.Find(&items).Error; err != nil {
		logger.GetLogger("quiver").Errorf("failed to query items: %s", err.Error())
		return nil, 0, fmt.Errorf("failed to query items: %w", err)
	}

	return items, total, nil
}

// DeleteItem 删除配置项
func (s *ItemService) DeleteItem(env, appName, clusterName, namespaceName, key string) error {
	// 检查 app, cluster, namespace
	ids, err := CheckACNKinDB(&env, &appName, &clusterName, &namespaceName, &key)
	if err != nil {
		return err
	}
	db := database.GetDB(env)

	return db.Where("namespace_id = ? AND k = ?", ids.NamespaceID, key).Delete(&models.Item{}).Error
}
