package model

import "time"

type Order struct {
	Id          int `gorm:"primary_key"`
	UserId      int
	CommodityId int
	Quantity    int
	TotalPrice  float64
	Address     string
	Status      string
	CreatedAt   time.Time
	UpdateAt    time.Time
}
