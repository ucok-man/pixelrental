package repo

import (
	"errors"
	"strings"

	"github.com/ucok-man/pixelrental/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CartService struct {
	db *gorm.DB
}

func (s *CartService) GetAll(userid int) ([]entity.Cart, error) {
	var carts []entity.Cart

	err := s.db.Preload("Game").Where("user_id = $1", userid).Find(&carts).Error
	if err != nil {
		return nil, err
	}

	if len(carts) == 0 {
		return make([]entity.Cart, 0), nil
	}

	return carts, nil
}

func (s *CartService) GetByGameID(userid, gameid int) (*entity.Cart, error) {
	var Cart *entity.Cart
	err := s.db.Where("user_id = $1 AND game_id = $2", userid, gameid).First(&Cart).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return Cart, nil
}

func (s *CartService) GetByID(userid, cartid int) (*entity.Cart, error) {
	var Cart *entity.Cart
	err := s.db.Where("user_id = $1 AND cart_id = $2", userid, cartid).First(&Cart).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return Cart, nil
}

func (s *CartService) Create(Cart *entity.Cart) error {
	err := s.db.Clauses(clause.Returning{}).Create(Cart).Error
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "violates foreign key constraint"):
			return ErrViolateForeignKey
		default:
			return err
		}
	}
	return nil
}

func (s *CartService) DeleteOne(cartid int) error {
	err := s.db.Delete(&entity.Cart{}, cartid).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *CartService) DeleteAll(userid int) error {
	err := s.db.Where("user_id = $1", userid).Delete(&entity.Cart{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *CartService) Update(cart *entity.Cart) error {
	err := s.db.
		Clauses(clause.Returning{}).
		Model(cart).
		Updates(entity.Cart{Quantity: cart.Quantity, SubTotal: cart.SubTotal}).
		Error
	if err != nil {
		return err
	}
	return nil
}
