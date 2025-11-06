package model

import "time"

// Order 订单模型
type Order struct{
	Id          int `gorm:"primary_key"`
	UserId      int
	CommodityId int
	Quantity    int
	TotalPrice  string
	Address     string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
