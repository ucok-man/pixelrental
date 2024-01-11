package entity

type OrderDetail struct {
	OrderDetailID int `gorm:"primaryKey"`
	OrderID       int
	GameID        int
	Quantity      int
	SubTotal      float64

	Game *Game `gorm:"foreignKey:GameID;References:GameID"`
}
