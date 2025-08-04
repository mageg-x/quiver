package services

import (
	"errors"
	"gorm.io/gorm"
	"quiver/database"
	"quiver/logger"
	"quiver/models"
)

// ItemService 命名空间服务
type ItemService struct{}

// NewItemService 创建命名空间服务实例
func NewItemService() *ItemService {
	return &ItemService{}
}

// CheckACN 检查 app, cluster, namespace
func CheckACN(env, appName, clusterName, namespaceName string) (*models.IDs, error) {
	db := database.GetDB(env)
	// 检查应用是否存在
	var app models.App
	err := db.Where("app_name = ?", appName).First(&app).Error
	if err != nil {
		logger.GetLogger("quiver").Errorf("app %s not found", appName)
		return nil, errors.New("app not found")
	}
	// 检查集群是否存在
	var cluster models.Cluster
	err = db.Where("app_name = ? AND cluster_name = ?", appName, clusterName).First(&cluster).Error
	if err != nil {
		logger.GetLogger("quiver").Errorf("cluster %s not found", clusterName)
		return nil, errors.New("cluster not found")
	}

	// 检查命名空间是否存在
	var namespace models.Namespace
	err = db.Where("namespace_name = ?", namespaceName).First(&namespace).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("namespace not found")
		}
		return nil, errors.New("namespace not found")
	}
	return &models.IDs{AppID: app.AppID, ClusterID: cluster.ClusterID, NamespaceID: namespace.NamespaceID}, err
}

// SetItem 设置配置项
func (s *ItemService) SetItem(env string, appName, clusterName, namespaceName, key, value string) error {
	// 检查 app, cluster, namespace
	ids, err := CheckACN(env, appName, clusterName, namespaceName)
	if err != nil {
		return err
	}

	db := database.GetDB(env)
	// 更新或创建配置项
	var item models.Item
	err = db.Where("namespace_id = ? AND k = ?", ids.NamespaceID, key).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新配置项
			item = models.Item{
				AppID:         ids.AppID,
				ClusterID:     ids.ClusterID,
				NamespaceID:   ids.NamespaceID,
				NamespaceName: namespaceName,
				K:             key,
				V:             value,
			}
			err = db.Create(&item).Error
		} else {
			logger.GetLogger("quiver").Errorf("failed to create item")
			return err
		}
	} else {
		// 更新现有配置项
		item.V = value
		// 更新现有配置项
		err = db.Model(&item).Update("V", value).Error
	}

	if err != nil {
		logger.GetLogger("quiver").Errorf("failed to update item")
		return err
	}

	return nil
}

// GetItem 获取单个配置项
func (s *ItemService) GetItem(env, appName, clusterName, namespaceName, key string) (string, error) {
	// 检查 app, cluster, namespace
	ids, err := CheckACN(env, appName, clusterName, namespaceName)
	if err != nil {
		return "", err
	}

	// 检查配置项是否存在
	db := database.GetDB(env)
	var item models.Item
	err = db.Where("namespace_id = ? AND k = ?", ids.NamespaceID, key).First(&item).Error
	if err != nil {
		logger.GetLogger("quiver").Errorf("item %s not found", key)
		return "", errors.New("item not found")
	}

	return item.V, nil
}

// DeleteItem 删除配置项
func (s *ItemService) DeleteItem(env, appName, clusterName, namespaceName, key string) error {
	// 检查 app, cluster, namespace
	ids, err := CheckACN(env, appName, clusterName, namespaceName)
	if err != nil {
		return err
	}
	db := database.GetDB(env)

	// 检查配置项是否存在
	var item models.Item
	err = db.Where("namespace_id = ? AND k = ?", ids.NamespaceID, key).First(&item).Error
	if err != nil {
		logger.GetLogger("quiver").Errorf("item %s not found", key)
		return errors.New("item not found")
	}

	return db.Where("namespace_id = ? AND k = ?", ids.NamespaceID, key).Delete(&models.Item{}).Error
}
