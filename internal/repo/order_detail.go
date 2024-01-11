package repo

import (
	"github.com/ucok-man/pixelrental/internal/entity"
	"gorm.io/gorm"
)

type OrderDetailService struct {
	db *gorm.DB
}

func (s *OrderDetailService) CreateBatch(orderDetails []*entity.OrderDetail) error {
	err := s.db.Create(orderDetails).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *OrderDetailService) GetAll(orderid int) ([]*entity.OrderDetail, error) {
	var ods []*entity.OrderDetail

	err := s.db.Preload("Game").Where("order_id = $1", orderid).Find(&ods).Error
	if err != nil {
		return nil, err
	}

	if len(ods) == 0 {
		return make([]*entity.OrderDetail, 0), nil
	}

	return ods, nil
}
