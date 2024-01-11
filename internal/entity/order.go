package entity

import "time"

const (
	ORDER_STATUS_WAITING   = "WAITING"
	ORDER_STATUS_SENDING   = "SENDING"
	ORDER_STATUS_FULLFILED = "FULLFILED"
)

type Order struct {
	OrderID    int `gorm:"primaryKey"`
	UserID     int
	Status     string
	TotalPrice float64
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`

	OrderDetails []*OrderDetail `gorm:"foreignKey:OrderID"`
	Payment      *Payment       `gorm:"foreignKey:OrderID"`
	User         *User          `gorm:"foreignKey:UserID;References:UserID"`
}
