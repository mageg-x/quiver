package models

import (
	"fmt"
	"time"
)

// Namespace 命名空间模型
type Namespace struct {
	NamespaceID   uint64    `json:"-" gorm:"column:id;primaryKey;autoIncrement"`
	NamespaceName string    `json:"namespace_name" gorm:"column:namespace_name;size:128;not null;index:idx_namespace_name"`
	Description   string    `json:"description" gorm:"column:description;size:1024"`
	AppID         uint64    `json:"-" gorm:"column:app_id;not null;index:idx_app_id"`
	AppName       string    `json:"app_name" gorm:"column:app_name;size:128;not null;index:idx_app_name"`
	ClusterID     uint64    `json:"-" gorm:"column:cluster_id;not null;index:idx_cluster_id"`
	ClusterName   string    `json:"cluster_name" gorm:"column:cluster_name;size:128;not null;index:idx_cluster_name"`
	Ver           uint64    `json:"-" gorm:"column:ver;not null;default:1"`
	CreateTime    time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`
	UpdateTime    time.Time `json:"update_time" gorm:"column:update_time;autoUpdateTime"`

	// 关联关系（可选加载）
	App     *App     `json:"-" gorm:"foreignKey:AppID;references:AppID"`
	Cluster *Cluster `json:"-" gorm:"foreignKey:ClusterID;references:ClusterID"`
}

func (n *Namespace) GetUpdateTime() time.Time {
	return n.UpdateTime
}

// TableName 指定表名
func (n *Namespace) TableName() string {
	return "namespace"
}

type NamespaceRelease struct {
	ID            uint64    `json:"-" gorm:"column:id;primaryKey;autoIncrement"`
	AppID         uint64    `json:"-" gorm:"column:app_id;not null;index:idx_app_id"`
	AppName       string    `json:"app_name" gorm:"column:app_name;size:128;not null"`
	ClusterID     uint64    `json:"-" gorm:"column:cluster_id;not null;index:idx_cluster_id"`
	ClusterName   string    `json:"cluster_name" gorm:"column:cluster_name;size:128;not null"`
	NamespaceID   uint64    `json:"-" gorm:"column:namespace_id;not null;index:idx_namespace_id"`
	NamespaceName string    `json:"namespace_name" gorm:"column:namespace_name;size:128;not null;index:idx_namespace_name"`
	ReleaseID     string    `json:"release_id" gorm:"column:release_id;size:64;not null;uniqueIndex:uk_release_id"`
	ReleaseName   string    `json:"release_name" gorm:"column:release_name;size:128;not null"`
	ReleaseTime   time.Time `json:"release_time" gorm:"column:create_time;autoCreateTime"`
	Operator      string    `json:"operator" gorm:"column:operator;size:64"`
	Comment       string    `json:"comment" gorm:"column:comment;type:varchar(1024)"`
	Config        []byte    `json:"-" gorm:"column:config;type:blob"`
	CreateTime    time.Time `json:"-" gorm:"column:create_time;autoCreateTime"`
	UpdateTime    time.Time `json:"-" gorm:"column:update_time;autoUpdateTime"`

	// 关联关系（可选加载）
	KvIDs   []uint64      `json:"-" gorm:"-"`
	Items   []ItemRelease `json:"-" gorm:"-"`
	App     App           `json:"-" gorm:"foreignKey:AppID;references:AppID"`
	Cluster Cluster       `json:"-" gorm:"foreignKey:ClusterID;references:AppID"`
}

func (nr *NamespaceRelease) GetID() uint64            { return nr.ID }
func (nr *NamespaceRelease) GetUpdateTime() time.Time { return nr.UpdateTime }
func (nr *NamespaceRelease) CacheKey(env string) string {
	return fmt.Sprintf("release:latest:%s:%s:%s:%s", env, nr.AppName, nr.ClusterName, nr.NamespaceName)
}

// TableName 指定表名
func (nr *NamespaceRelease) TableName() string {
	return "namespace_release"
}
