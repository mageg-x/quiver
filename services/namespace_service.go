package services

import (
	"errors"
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"quiver/database"
	"quiver/logger"
	"quiver/models"
	"quiver/utils"
)

// NamespaceService 命名空间服务
type NamespaceService struct{}

// NewNamespaceService 创建命名空间服务实例
func NewNamespaceService() *NamespaceService {
	return &NamespaceService{}
}

// CreateNamespace 创建命名空间
func (s *NamespaceService) CreateNamespace(env string, namespace *models.Namespace) error {
	ids, err := CheckACNKinDB(&env, &namespace.AppName, &namespace.ClusterName, nil, nil)
	if err != nil {
		return err
	}

	namespace.AppID = ids.AppID
	namespace.ClusterID = ids.ClusterID

	db := database.GetDB(env)
	// 检查命名空间是否已存在
	var existingNamespace models.Namespace
	err = db.Where("app_name =? AND cluster_name = ? AND namespace_name = ?", namespace.AppName, namespace.ClusterName, namespace.NamespaceName).First(&existingNamespace).Error
	if err == nil {
		logger.GetLogger("quiver").Errorf("namespace %s already exists in this cluster", namespace.NamespaceName)
		return errors.New("namespace already exists")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.GetLogger("quiver").Errorf("namespace create %s failed %s", namespace.NamespaceName, err.Error())
		return err
	}

	// 创建命名空间
	return db.Create(namespace).Error
}

// ListNamespace 获取集群下的所有命名空间
func (s *NamespaceService) ListNamespace(env string, appName, clusterName string, page, size int) ([]models.Namespace, int64, error) {
	ids, err := CheckACNKinDB(&env, &appName, &clusterName, nil, nil)
	if err != nil {
		return nil, 0, err
	}

	db := database.GetDB(env)

	var namespaces []models.Namespace
	var total int64

	// 构建查询条件
	query := db.Where("cluster_id = ?", ids.ClusterID).Order("update_time DESC,id DESC")

	// 获取总数（用于分页）
	if err := query.Model(&models.Namespace{}).Count(&total).Error; err != nil {
		logger.GetLogger("quiver").Errorf("failed to count namespaces num : %s for cluster: %s", err.Error(), clusterName)
		return nil, 0, fmt.Errorf("failed to count namespaces: %w", err)
	}

	offset := (page - 1) * size

	// 执行分页查询
	if err := query.Offset(offset).Limit(size).Find(&namespaces).Error; err != nil {
		logger.GetLogger("quiver").Errorf("failed to query namespaces: %s", err.Error())
		return nil, 0, fmt.Errorf("failed to query namespaces: %w", err)
	}

	return namespaces, total, nil
}

// GetNamespace 获取特定命名空间
func (s *NamespaceService) GetNamespace(env string, appName, clusterName, namespaceName string) (*models.Namespace, error) {
	ids, err := CheckACNKinDB(&env, &appName, &clusterName, &namespaceName, nil)
	if err != nil {
		return nil, err
	}

	db := database.GetDB(env)

	var namespace models.Namespace
	err = db.Where("id = ?", ids.NamespaceID).First(&namespace).Error
	if err != nil {
		logger.GetLogger("quiver").Errorf("namespace %s not found",
			appName+"/"+clusterName+"/"+namespaceName)
		return nil, errors.New("namespace not found")
	}

	return &namespace, nil
}

// DeleteNamespace 删除命名空间
func (s *NamespaceService) DeleteNamespace(env string, appName, clusterName, namespaceName string) error {
	ids, err := CheckACNKinDB(&env, &appName, &clusterName, &namespaceName, nil)
	if err != nil {
		return err
	}

	db := database.GetDB(env)

	return db.Transaction(func(tx *gorm.DB) error {
		// 更新  Item 记录的 deleted 字段
		err := BulkMarkDeleted[*models.Item](tx, map[string]interface{}{"namespace_id": ids.NamespaceID})
		if err != nil {
			logger.GetLogger("quiver").Errorf("update items deleted for namespace %s failed %s", namespaceName, err)
			return err
		}

		// 更新与  ItemRelease 记录的 deleted 字段
		err = BulkMarkDeleted[*models.ItemRelease](tx, map[string]interface{}{"namespace_id": ids.NamespaceID})
		if err != nil {
			logger.GetLogger("quiver").Errorf("update items release deleted for namespace %s failed %s", namespaceName, err)
			return err
		}

		// 删除命名空间
		result := tx.Where("id = ?", ids.NamespaceID).Delete(&models.Namespace{})
		if result.Error != nil {
			logger.GetLogger("quiver").Errorf("delete namespace %s failed %s", namespaceName, result.Error)
			return result.Error
		}

		if result.RowsAffected == 0 {
			logger.GetLogger("quiver").Errorf("namespace delete %s failed %s", namespaceName, "namespace not found")
			return errors.New("namespace not found")
		}

		return nil
	})
}

func (s *NamespaceService) DiscardDraft(env, appName, clusterName, namespaceName string) error {
	ids, err := CheckACNKinDB(&env, &appName, &clusterName, &namespaceName, nil)
	if err != nil {
		return err
	}

	db := database.GetDB(env)

	err = BulkDeleted[*models.Item](db, map[string]interface{}{"namespace_id": ids.NamespaceID})
	if err != nil {
		logger.GetLogger("quiver").Errorf("BulkDelete %s item table failed %s", namespaceName, err)
	}

	err = BulkDeleted[*models.ItemRelease](db, map[string]interface{}{"namespace_id": ids.NamespaceID})
	if err != nil {
		logger.GetLogger("quiver").Errorf("BulkDelete %s item release table failed %s", namespaceName, err)
	}

	// 1、获取 当前namespace 中所有配置项目
	var items []models.Item
	if err := db.Where("namespace_id = ? AND is_deleted =0", ids.NamespaceID).Find(&items).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		logger.GetLogger("quiver").Errorf("get items failed: %v", err)
		return err
	}

	if len(items) == 0 {
		return nil
	}

	var draftK, updateK []uint32
	drafMap := make(map[uint32]models.Item, len(items))
	for _, item := range items {
		kvID := item.KVId
		hashK := uint32(kvID & 0xFFFFFFFF) // 取低 32 位
		draftK = append(draftK, hashK)
		if item.IsReleased == 0 {
			updateK = append(updateK, hashK)
		}
		drafMap[hashK] = item
	}

	// 2. 获取最近一次的 release
	var latestRelease models.NamespaceRelease
	if err := db.Where("namespace_id = ?", ids.NamespaceID).
		Order("id DESC").
		First(&latestRelease).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.GetLogger("quiver").Infof("no releases found for %s/%s/%s/%s", env, appName, clusterName, namespaceName)
		} else {
			// 处理其他可能的错误
			logger.GetLogger("quiver").Errorf("error while fetching latest release: %v", err)
			return err
		}
	}

	if latestRelease.ReleaseID == "" {
		// 没有发布过任何版本，直接删除当前草稿区域
		result := db.Where("namespace_id = ?", ids.NamespaceID).Delete(&models.Item{})

		if result.Error != nil {
			logger.GetLogger("quiver").Errorf("delete all items from item table failed: %v", result.Error)
			return result.Error
		}

		logger.GetLogger("quiver").Infof("deleted %d items from item table", result.RowsAffected)
		return nil
	} else {
		var latestKvIDs []uint64
		if len(latestRelease.Config) > 0 {
			if err := msgpack.Unmarshal(latestRelease.Config, &latestKvIDs); err != nil {
				logger.GetLogger("quiver").Errorf("failed to unmarshal %s config", latestRelease.ReleaseID)
				return errors.New("failed to get release data")
			}
		}
		var latestK []uint32
		latestKK := make(map[uint32]uint64, len(latestKvIDs))
		for _, kvID := range latestKvIDs {
			hashK := uint32(kvID & 0xFFFFFFFF) // 取低 32 位
			latestK = append(latestK, hashK)
			latestKK[hashK] = kvID
		}

		// 计算变化部分
		delK, _, addK := utils.Diff32(draftK, latestK)
		_, updateK, _ = utils.Diff32(updateK, latestK)
		// 把delK 转成 delId，再批量删除
		var delId []uint64
		for _, k := range delK {
			if m, ok := drafMap[k]; ok {
				delId = append(delId, m.ID)
			}
		}

		// 把 addK 转成 addKvIDs
		var addKvIDs []uint64
		for _, k := range addK {
			if m, ok := latestKK[k]; ok {
				addKvIDs = append(addKvIDs, m)
			}
		}

		// 把updateK 转成 updateKvIDs
		var updateKvIDs []uint64
		for _, k := range updateK {
			if m, ok := latestKK[k]; ok {
				updateKvIDs = append(updateKvIDs, m)
			}
		}
		// 4. 开启事务
		tx := db.Begin()
		if tx.Error != nil {
			logger.GetLogger("quiver").Errorf("begin transaction failed: %v", tx.Error)
			return fmt.Errorf("failed to begin transaction: %w", tx.Error)
		}

		defer func() {
			if tx != nil {
				tx.Rollback()
			}
		}()

		// 5、从item release 表中 批量读取 addKvIDs + updateKvIDs
		var allItemReleases []models.ItemRelease
		allKvIDs := utils.Merge(addKvIDs, updateKvIDs)
		if len(allKvIDs) > 0 {
			if err := tx.Where("kv_id IN ?", allKvIDs).Find(&allItemReleases).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("item not found: %w", err)
				}
				logger.GetLogger("quiver").Errorf("query items failed: %v", err)
				return err
			}
		}

		var allItems []models.Item
		for _, itemRelease := range allItemReleases {
			item := models.Item{
				AppID:       itemRelease.AppID,
				ClusterID:   itemRelease.ClusterID,
				NamespaceID: itemRelease.NamespaceID,
				K:           itemRelease.K,
				V:           itemRelease.V,
				KVId:        itemRelease.KvID,
				IsReleased:  1,
				IsDeleted:   0,
			}
			allItems = append(allItems, item)
		}
		// 6、批量从item 表中删除 delId
		if len(delId) > 0 {
			// 假设你已经有 namespaceID 变量
			result := tx.Where("namespace_id = ? AND id IN ?", ids.NamespaceID, delId).Delete(&models.Item{})

			if result.Error != nil {
				logger.GetLogger("quiver").Errorf("batch delete items from item table failed: %v", result.Error)
				return result.Error
			}

			logger.GetLogger("quiver").Infof("batch deleted %d items from item table", result.RowsAffected)
		}

		// 7、批量把 allItems 中的 内容插入 到 item 表中，如果存在则更新
		if len(allItems) > 0 {
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{Name: "namespace_id"},
					{Name: "k"},
				},
				DoUpdates: clause.AssignmentColumns([]string{"v", "kv_id", "is_released", "is_deleted"}),
			}).Create(&allItems).Error; err != nil {
				logger.GetLogger("quiver").Errorf("batch upsert items failed: %v", err)
				return err
			}
		}
		// 8. 提交事务
		if err := tx.Commit().Error; err != nil {
			logger.GetLogger("quiver").Errorf("transaction commit failed: %v", err)
			return err
		}
		tx = nil
	}

	return err
}
