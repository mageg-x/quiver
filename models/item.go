package models

import (
	"time"
)

type HasID interface {
	GetID() uint64
}

// Item 配置项模型
type Item struct {
	ID          uint64 `json:"-" gorm:"column:id;primaryKey;autoIncrement"`
	AppID       uint64 `json:"-" gorm:"column:app_id;not null;index:idx_app_id"`
	ClusterID   uint64 `json:"-" gorm:"column:cluster_id;not null;index:idx_cluster_id"`
	NamespaceID uint64 `json:"-" gorm:"column:namespace_id;not null;index:idx_namespace_id"`
	K           string `json:"key" gorm:"column:k;size:255;not null;index:idx_k"`
	V           string `json:"value" gorm:"column:v;type:text"`
	KVId        uint64 `json:"-" gorm:"column:kv_id;not null;index:idx_kv_id"` // MurmurHash64(k,v)

	CreateTime time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time;autoUpdateTime"`
	IsDeleted  uint8     `json:"-" gorm:"column:is_deleted;not null;default:0"`            // 逻辑删除
	IsReleased uint8     `json:"is_released" gorm:"column:is_released;not null;default:0"` // 是否已发布

	// 关联字段（可选）
	Namespace *Namespace `json:"namespace,omitempty" gorm:"foreignKey:NamespaceID;references:NamespaceID"`
}

func (i *Item) GetID() uint64 {
	return i.ID
}

func (i *Item) GetUpdateTime() time.Time {
	return i.UpdateTime
}

// TableName 指定表名
func (i *Item) TableName() string {
	return "item"
}

// ItemRelease 发布配置项模型
type ItemRelease struct {
	ID          uint64 `json:"-" gorm:"column:id;primaryKey;autoIncrement"`
	AppID       uint64 `json:"-" gorm:"column:app_id;not null"`
	ClusterID   uint64 `json:"-" gorm:"column:cluster_id;not null"`
	NamespaceID uint64 `json:"-" gorm:"column:namespace_id;not null"`
	K           string `json:"key" gorm:"column:k;size:255;not null"`
	V           string `json:"value" gorm:"column:v;type:text"`
	KvID        uint64 `json:"kv_id" gorm:"column:kv_id;not null;index:idx_kv_id"`     // 使用 uint64 代替 BIGINT
	IsDeleted   uint8  `json:"is_deleted" gorm:"column:is_deleted;not null;default:0"` // 逻辑删除

	CreateTime time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time;autoCreateTime"`
}

func (ir *ItemRelease) GetUpdateTime() time.Time {
	return ir.UpdateTime
}

func (ir *ItemRelease) GetID() uint64 {
	return ir.ID
}

// TableName 指定表名
func (ir *ItemRelease) TableName() string {
	return "item_release"
}
