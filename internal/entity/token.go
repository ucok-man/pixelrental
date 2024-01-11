package entity

import (
	"time"
)

const (
	ScopeActivation    = "activation"
	ScopePasswordReset = "password-reset"
	ScopeInvoiceTopUp  = "invoice-topup"
	ScopeInvoiceOrder  = "invoice-order"
)

type Token struct {
	TokenHash []byte `gorm:"primaryKey"`
	Plaintext string `gorm:"-:all"`
	UserID    int
	Scope     string
	Expiry    time.Time

	// association
	User *User `gorm:"foreignKey:UserID;References:UserID"`
}
