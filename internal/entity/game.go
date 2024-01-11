package entity

import (
	"time"

	"github.com/lib/pq"
)

type Game struct {
	GameID      int `gorm:"primaryKey"`
	Title       string
	Description string
	Price       float64
	Year        int32
	Genres      pq.StringArray `gorm:"type:text[]"`
	Stock       int32
	Version     int32
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
