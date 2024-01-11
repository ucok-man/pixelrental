package repo

import (
	"errors"
	"strings"

	"github.com/ucok-man/pixelrental/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CalculatePriceTotal(carts []entity.Cart) float64 {
	var total float64
	for _, cart := range carts {
		total += cart.SubTotal
	}
	return total
}

type OrderService struct {
	db *gorm.DB
}

func (s *OrderService) GetAllByUserID(userid int) ([]*entity.Order, error) {
	var orders []*entity.Order
	err := s.db.
		Preload("User").
		Preload("Payment").
		Preload("OrderDetails.Game").
		Preload("OrderDetails").
		Where("user_id = $1", userid).
		Find(&orders).Error

	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return make([]*entity.Order, 0), nil
	}
	return orders, nil
}

func (s *OrderService) Create(Order *entity.Order) error {
	err := s.db.Create(Order).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderService) GetByID(orderid int) (*entity.Order, error) {
	Order := entity.Order{}

	err := s.db.Where("order_id = $1", orderid).First(&Order).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &Order, nil
}

// NOTE: Implement versioning to prevent race update
func (s *OrderService) UpdateStatus(Order *entity.Order) error {
	// versionbefore := user.Version
	err := s.db.
		Clauses(clause.Returning{}).
		Model(Order).
		// Where("version = $11", versionbefore).
		Updates(&entity.Order{
			Status: Order.Status,
		}).Error
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "duplicate key value violates unique constraint"):
			return ErrDuplicateRecord
		case errors.Is(err, gorm.ErrRecordNotFound):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}
