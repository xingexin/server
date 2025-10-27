package model

import "time"

type User struct {
	Uid       int `gorm:"primaryKey"`
	Account   string
	Password  string
	Name      string
	CreatedAt time.Time
}
