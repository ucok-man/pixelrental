package entity

import "time"

type Cart struct {
	CartID    int `gorm:"primaryKey"`
	UserID    int
	GameID    int
	Quantity  int
	SubTotal  float64
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Game *Game `gorm:"foreignKey:GameID;References:GameID"`
}
