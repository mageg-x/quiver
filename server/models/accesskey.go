package models

import (
	"fmt"
	"time"
)

type AccessKey struct {
	ID         uint64    `json:"-" gorm:"column:id;primaryKey;autoIncrement"`
	UserID     uint64    `json:"user_id" gorm:"column:user_id;not null"`
	AccessKey  string    `json:"access_key" gorm:"column:access_key;size:64;not null;uniqueIndex:uk_access_key"`
	SecretKey  string    `json:"secret_key,omitempty" gorm:"column:secret_key;size:128;not null"` // 建议存储加密后的值
	CreateTime time.Time `json:"-" gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time `json:"-" gorm:"column:update_time;autoUpdateTime"`
}

func (ak *AccessKey) GetID() uint64            { return ak.ID }
func (ak *AccessKey) GetUpdateTime() time.Time { return ak.UpdateTime }
func (ak *AccessKey) CacheKey(env string) string {
	return fmt.Sprintf("accesskey:%s:%d:%s", env, ak.UserID, ak.AccessKey)
}

// TableName 指定表名
func (ak *AccessKey) TableName() string {
	return "accesskey"
}
