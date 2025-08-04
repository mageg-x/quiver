package models

import (
	"time"
)

// Item 配置项模型
type Item struct {
	ItemID        uint64    `json:"-" gorm:"column:id;primaryKey;autoIncrement"`
	AppID         uint64    `json:"-" gorm:"column:app_id;not null"`
	ClusterID     uint64    `json:"-" gorm:"column:cluster_id;not null"`
	NamespaceID   uint64    `json:"-" gorm:"column:namespace_id;not null"`
	K             string    `json:"k" gorm:"column:k;size:255;not null;index:idx_k"`
	V             string    `json:"v" gorm:"column:v;type:text"`
	NamespaceName string    `json:"namespace_name" gorm:"column:namespace_name;size:128;not null;index:idx_namespace_name"`
	Ver           uint64    `json:"-" gorm:"column:ver;not null;default:1"`
	CreateTime    time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`
	UpdateTime    time.Time `json:"update_time" gorm:"column:update_time;autoUpdateTime"`
	Deleted       uint8     `json:"-" gorm:"column:deleted;not null;default:0"` // 逻辑删除字段

	Namespace *Namespace `json:"namespace,omitempty" gorm:"foreignKey:NamespaceID;references:NamespaceID"`
}

// TableName 指定表名
func (Item) TableName() string {
	return "item"
}

// ItemRelease 发布配置项模型
type ItemRelease struct {
	ID          uint64 `json:"-" gorm:"column:id;primaryKey;autoIncrement"`
	AppID       uint64 `json:"-" gorm:"column:app_id;not null"`
	ClusterID   uint64 `json:"-" gorm:"column:cluster_id;not null"`
	NamespaceID uint64 `json:"-" gorm:"column:namespace_id;not null"`
	ReleaseID   string `json:"release_id" gorm:"column:release_id;size:64;not null;index:idx_release_id"`
	K           string `json:"k" gorm:"column:k;size:255;not null"`
	V           string `json:"v" gorm:"column:v;type:text"`
	KvID        uint64 `json:"kv_id" gorm:"column:kv_id;not null;index:idx_kv_id"` // 使用 uint64 代替 BIGINT
	Deleted     uint8  `json:"-" gorm:"column:deleted;not null;default:0"`         // 逻辑删除字段

	CreateTime time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time;autoCreateTime"`
}

// TableName 指定表名
func (ItemRelease) TableName() string {
	return "item_release"
}
