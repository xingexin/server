package model

import "time"

// Commodity 商品模型
type Commodity struct {
	ID        int     `gorm:"primaryKey" json:"id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
	Status    bool
	CreatedAt time.Time
	UpdateAt  time.Time
}
