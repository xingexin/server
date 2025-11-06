package model

import "time"

// User 用户模型
type User struct {
	Uid       int `gorm:"primaryKey"`
	Account   string
	Password  string
	Name      string
	CreatedAt time.Time
}
