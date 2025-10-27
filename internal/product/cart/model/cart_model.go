package model

import "time"

type Cart struct {
	Id          int `gorm:"primaryKey"`
	UserId      int
	CommodityId int
	Quantity    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
