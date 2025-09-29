package model

import "time"

type Commodity struct {
	ID        int     `gorm:"primaryKey"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Stock     int
	Status    bool
	CreatedAt time.Time
	UpdateAt  time.Time
}
