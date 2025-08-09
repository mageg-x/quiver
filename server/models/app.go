package models

import (
	"time"
)

type IDs struct {
	AppID       uint64
	ClusterID   uint64
	NamespaceID uint64
	ItemID      uint64
}

// App 应用模型
type App struct {
	AppID       uint64    `json:"app_id" gorm:"column:id;primaryKey;autoIncrement"`
	AppName     string    `json:"app_name" gorm:"column:app_name;size:128;not null;uniqueIndex:uk_app_name"`
	Description string    `json:"description" gorm:"column:description;size:1024"`
	Ver         uint64    `json:"-" gorm:"column:ver;not null;default:1"` // CAS 乐观锁版本号
	CreateTime  time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`
	UpdateTime  time.Time `json:"update_time" gorm:"column:update_time;autoUpdateTime"`
}

func (a *App) GetUpdateTime() time.Time {
	return a.UpdateTime
}

// TableName 指定表名
func (a *App) TableName() string {
	return "app"
}
