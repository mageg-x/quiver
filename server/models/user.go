package models

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User 用户模型
type User struct {
	UserID     uint64    `json:"user_id" gorm:"column:id;primaryKey;autoIncrement"`
	UserName   string    `json:"user_name" gorm:"column:user_name;uniqueIndex;size:128;not null"`
	Password   string    `json:"password,omitempty" gorm:"column:password;size:256;not null"`
	Email      string    `json:"email" gorm:"column:email;size:128"`
	Phone      string    `json:"phone" gorm:"column:phone;size:32"`
	CreateTime time.Time `json:"-" gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time `json:"-" gorm:"column:update_time;autoUpdateTime"`
}

// SetPassword 设置密码（加密）
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// TableName 指定表名
func (u *User) TableName() string {
	return "user"
}
