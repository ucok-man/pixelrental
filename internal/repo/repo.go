package repo

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrRecordNotFound    = errors.New("no record found")
	ErrDuplicateRecord   = errors.New("record duplicate on unique constraint")
	ErrViolateForeignKey = errors.New("record violate foreign key constraint")
	ErrEditConflict      = errors.New("record edit conflict when updating")
)

type Services struct {
	User        UserService
	Game        GameServices
	Token       TokenService
	Cart        CartService
	Order       OrderService
	OrderDetail OrderDetailService
	Payment     PaymentService
}

func New(db *gorm.DB) *Services {
	return &Services{
		User:        UserService{db: db},
		Game:        GameServices{db: db},
		Token:       TokenService{db: db},
		Cart:        CartService{db: db},
		Order:       OrderService{db: db},
		OrderDetail: OrderDetailService{db: db},
		Payment:     PaymentService{db: db},
	}
}
