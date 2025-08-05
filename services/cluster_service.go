package services

import (
	"errors"
	"fmt"
	"quiver/database"
	"quiver/logger"
	"quiver/models"

	"gorm.io/gorm"
)

// ClusterService 集群服务
type ClusterService struct{}

// NewClusterService 创建集群服务实例
func NewClusterService() *ClusterService {
	return &ClusterService{}
}

// CreateCluster 创建集群
func (s *ClusterService) CreateCluster(env string, cluster *models.Cluster) error {
	ids, err := CheckACNKinDB(&env, &cluster.AppName, nil, nil, nil)
	if err != nil {
		return err
	}

	db := database.GetDB(env)

	cluster.AppID = ids.AppID

	// 检查集群名称是否已存在
	var existingCluster models.Cluster
	err = db.Where("app_name = ? AND cluster_name = ?", cluster.AppName, cluster.ClusterName).First(&existingCluster).Error
	if err == nil {
		logger.GetLogger("quiver").Errorf("cluster %s already exists in this app", cluster.ClusterName)
		return errors.New("cluster already exists in this app")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.GetLogger("quiver").Errorf("cluster create %s failed %s", cluster.AppName, err.Error())
		return err
	}

	// 创建集群
	return db.Create(cluster).Error
}

// ListCluster 获取应用下的所有集群
func (s *ClusterService) ListCluster(env string, appName string, page, size int) ([]models.Cluster, int64, error) {
	ids, err := CheckACNKinDB(&env, &appName, nil, nil, nil)
	if err != nil {
		return nil, 0, err
	}

	// 获取对应环境的 DB 实例
	db := database.GetDB(env)

	var clusters []models.Cluster
	var total int64

	// 构建查询条件
	query := db.Where("app_id = ?", ids.AppID).Order("update_time DESC,id DESC")

	// 获取总数（用于分页）
	if err := query.Model(&models.Cluster{}).Count(&total).Error; err != nil {
		logger.GetLogger("quiver").Errorf("failed to count clusters for app %s: %v", appName, err)
		return nil, 0, fmt.Errorf("failed to count clusters: %w", err)
	}

	// 再查询分页数据
	offset := (page - 1) * size
	err = db.Where("app_name = ?", appName).
		Offset(offset).
		Limit(size).
		Order("create_time DESC"). // 可选：按创建时间倒序
		Find(&clusters).Error

	if err != nil {
		logger.GetLogger("quiver").Errorf("failed to list clusters for app %s: %v", appName, err)
		return nil, 0, err
	}

	return clusters, total, nil
}

// GetCluster 获取特定集群
func (s *ClusterService) GetCluster(env string, appName, clusterName string) (*models.Cluster, error) {
	ids, err := CheckACNKinDB(&env, &appName, &clusterName, nil, nil)
	if err != nil {
		return nil, err
	}

	db := database.GetDB(env)

	var cluster models.Cluster
	err = db.Where("id = ?", ids.ClusterID).First(&cluster).Error
	if err != nil {
		logger.GetLogger("quiver").Errorf("cluster get %s failed %s", appName, err.Error())
		return nil, err
	}

	return &cluster, nil
}

// DeleteCluster 删除集群
func (s *ClusterService) DeleteCluster(env string, appName, clusterName string) error {
	ids, err := CheckACNKinDB(&env, &appName, &clusterName, nil, nil)
	if err != nil {
		return err
	}
	db := database.GetDB(env)

	return db.Transaction(func(tx *gorm.DB) error {
		// 更新  Item 记录的 deleted 字段
		err := BulkMarkDeleted[models.Item](tx, map[string]interface{}{"cluster_id": ids.ClusterID})
		if err != nil {
			logger.GetLogger("quiver").Errorf("update items deleted for cluster %s failed %s", clusterName, err)
			return err
		}

		// 更新与  ItemRelease 记录的 deleted 字段
		err = BulkMarkDeleted[models.ItemRelease](tx, map[string]interface{}{"cluster_id": ids.ClusterID})
		if err != nil {
			logger.GetLogger("quiver").Errorf("update items release deleted for cluster %s failed %s", clusterName, err)
			return err
		}

		// 删除集群
		result := tx.Where("id = ?", ids.ClusterID).Delete(&models.Cluster{})
		if result.Error != nil {
			logger.GetLogger("quiver").Errorf("cluster delete %s failed %s", appName, result.Error.Error())
			return result.Error
		}

		if result.RowsAffected == 0 {
			logger.GetLogger("quiver").Errorf("cluster delete %s failed %s", appName, "cluster not found")
			return errors.New("cluster not found")
		}

		return nil
	})
}
