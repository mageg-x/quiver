package models

import (
	"time"
)

// Cluster 集群模型
type Cluster struct {
	ClusterID   uint64    `json:"cluster_id" gorm:"column:id;primaryKey;autoIncrement"`
	AppID       uint64    `json:"app_id" gorm:"column:app_id;not null;index:idx_app_id"`                            // 外键字段
	AppName     string    `json:"app_name" gorm:"column:app_name;size:128;not null;index:idx_app_name"`             // 用于唯一约束
	ClusterName string    `json:"cluster_name" gorm:"column:cluster_name;size:128;not null;index:idx_cluster_name"` // 同上
	Description string    `json:"description" gorm:"column:description;size:1024"`
	Ver         uint64    `json:"-" gorm:"column:ver;not null;default:1"` // 乐观锁版本号
	CreateTime  time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`
	UpdateTime  time.Time `json:"update_time" gorm:"column:update_time;autoUpdateTime"`

	// 关联关系：可选，按需加载
	App *App `json:"-" gorm:"foreignKey:AppID;references:AppID"`
}

func (c *Cluster) GetUpdateTime() time.Time {
	return c.UpdateTime
}

// TableName 指定表名
func (c *Cluster) TableName() string {
	return "cluster"
}
