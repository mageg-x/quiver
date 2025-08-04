package services

import (
	"errors"
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
	db := database.GetDB(env)

	// 检查应用是否存在
	var app models.App
	err := db.Where("app_name = ?", cluster.AppName).First(&app).Error
	if err != nil {
		logger.GetLogger("quiver").Errorf("app %s not found", app.AppName)
		return errors.New("app not found")
	}

	cluster.AppID = app.AppID

	// 检查集群名称是否已存在
	var existingCluster models.Cluster
	err = db.Where("app_name = ? AND cluster_name = ?", cluster.AppName, cluster.ClusterName).First(&existingCluster).Error
	if err == nil {
		logger.GetLogger("quiver").Errorf("cluster %s already exists in this app", cluster.ClusterName)
		return errors.New("cluster already exists in this app")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.GetLogger("quiver").Errorf("cluster create %s failed %s", app.AppName, err.Error())
		return err
	}

	// 创建集群
	return db.Create(cluster).Error
}

// ListCluster 获取应用下的所有集群
func (s *ClusterService) ListCluster(env string, appName string, page, size int) ([]models.Cluster, int64, error) {
	// 获取对应环境的 DB 实例
	db := database.GetDB(env)

	var clusters []models.Cluster
	var total int64

	// 检查应用是否存在
	var app models.App
	err := db.Where("app_name = ?", appName).First(&app).Error
	if err != nil {
		logger.GetLogger("quiver").Errorf("app %s not found", app.AppName)
		return nil, 0, errors.New("app not found")
	}

	// 先查询总数
	if err := db.Model(&models.Cluster{}).Where("app_name = ?", appName).Count(&total).Error; err != nil {
		logger.GetLogger("quiver").Errorf("failed to count clusters for app %s: %v", appName, err)
		return nil, 0, err
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
	db := database.GetDB(env)

	var cluster models.Cluster
	err := db.Where("app_name = ? AND cluster_name = ?", appName, clusterName).First(&cluster).Error
	if err != nil {
		logger.GetLogger("quiver").Errorf("cluster get %s failed %s", appName, err.Error())
		return nil, err
	}

	return &cluster, nil
}

// DeleteCluster 删除集群
func (s *ClusterService) DeleteCluster(env string, appName, clusterName string) error {
	db := database.GetDB(env)

	// 先检查集群是否存在
	var cluster models.Cluster
	err := db.Where("app_name = ? AND cluster_name = ?", appName, clusterName).First(&cluster).Error
	if err != nil {
		logger.GetLogger("quiver").Errorf("cluster delete %s failed %s", appName, err.Error())
		return errors.New("cluster not found")
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// 更新  Item 记录的 deleted 字段
		itemResult := tx.Model(&models.Item{}).Where("cluster_id = ?", cluster.ClusterID).Update("deleted", 1)
		if itemResult.Error != nil {
			logger.GetLogger("quiver").Errorf("update items deleted for cluster %s failed %s", clusterName, itemResult.Error)
			return itemResult.Error
		}
		if itemResult.RowsAffected == 0 {
			logger.GetLogger("quiver").Warnf("no items found for cluster %s", cluster.ClusterName)
		}

		// 更新与  ItemRelease 记录的 deleted 字段
		itemReleaseResult := tx.Model(&models.ItemRelease{}).Where("cluster_id = ?", cluster.ClusterID).Update("deleted", 1)
		if itemReleaseResult.Error != nil {
			logger.GetLogger("quiver").Errorf("update items release deleted for cluster %s failed %s", clusterName, itemReleaseResult.Error)
			return itemReleaseResult.Error
		}
		if itemReleaseResult.RowsAffected == 0 {
			logger.GetLogger("quiver").Warnf("no items releases found for cluster %s", clusterName)
		}

		// 删除集群
		result := tx.Where("app_name = ? AND cluster_name = ?", appName, clusterName).Delete(&models.Cluster{})
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
