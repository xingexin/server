package model

import "time"

type Commodity struct {
	ID        int `gorm:"primaryKey"`
	name      string
	price     float64
	stock     int
	status    bool
	createdAt time.Time
	updateAt  time.Time
}
