package services

import (
	"errors"
	"fmt"
	"quiver/database"
	"quiver/logger"
	"quiver/models"

	"gorm.io/gorm"
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
		err := BulkMarkDeleted[models.Item](tx, map[string]interface{}{"namespace_id": ids.NamespaceID})
		if err != nil {
			logger.GetLogger("quiver").Errorf("update items deleted for namespace %s failed %s", namespaceName, err)
			return err
		}

		// 更新与  ItemRelease 记录的 deleted 字段
		err = BulkMarkDeleted[models.ItemRelease](tx, map[string]interface{}{"namespace_id": ids.NamespaceID})
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
