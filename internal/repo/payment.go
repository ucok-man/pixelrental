package repo

import (
	"errors"
	"strings"

	"github.com/ucok-man/pixelrental/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PaymentService struct {
	db *gorm.DB
}

func (s *PaymentService) Create(Payment *entity.Payment) error {
	err := s.db.Create(Payment).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *PaymentService) GetByorderID(orderid int) (*entity.Payment, error) {
	var payment *entity.Payment

	err := s.db.Where("order_id = $1", orderid).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return payment, nil
}

// NOTE: Implement versioning to prevent race update
func (s *PaymentService) UpdateStatus(payment *entity.Payment) error {
	// versionbefore := user.Version
	err := s.db.
		Clauses(clause.Returning{}).
		Model(payment).
		// Where("version = $11", versionbefore).
		Updates(&entity.Payment{
			Status: payment.Status,
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
