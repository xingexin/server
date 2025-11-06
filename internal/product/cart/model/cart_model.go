package model

import "time"

// Cart 购物车模型
type Cart struct {
	Id          int `gorm:"primaryKey"`
	UserId      int
	CommodityId int
	Quantity    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
