package models

import "time"

type Permission struct {
	ID           uint64    `json:"permission_id" gorm:"column:id;primaryKey;autoIncrement"`
	UserID       uint64    `json:"user_id" gorm:"column:user_id;not null"`
	ResourceType string    `json:"resource_type" gorm:"column:resource_type;type:enum('APP', 'CLUSTER', 'NAMESPACE');not null"`
	ResourceID   uint64    `json:"resource_id" gorm:"column:resource_id;not null"`
	ResourceName string    `json:"resource_name" gorm:"column:resource_name;size:128;not null"`
	Action       string    `json:"action" gorm:"column:action;not null"`
	CreateTime   time.Time `json:"-" gorm:"column:create_time;autoCreateTime"`
	UpdateTime   time.Time `json:"-" gorm:"column:update_time;autoUpdateTime"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permission"
}
