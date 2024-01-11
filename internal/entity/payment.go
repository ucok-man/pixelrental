package entity

import "time"

const (
	PAYMENT_STATUS_PENDING = "PENDING"
	PAYMENT_STATUS_EXPIRED = "EXPIRED"
	PAYMENT_STATUS_PAID    = "PAID"
)

type Payment struct {
	PaymentID     int `gorm:"primaryKey"`
	OrderID       int
	Status        string
	PaymentMethod string
	InvoiceID     string
	InvoiceUrl    string
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}
