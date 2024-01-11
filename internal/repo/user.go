package repo

import (
	"crypto/sha256"
	"errors"
	"strings"
	"time"

	"github.com/ucok-man/pixelrental/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserService struct {
	db *gorm.DB
}

func (s *UserService) GetByEmail(email string) (*entity.User, error) {
	User := entity.User{}

	err := s.db.Where("email = $1", email).First(&User).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &User, nil
}

func (s *UserService) GetByID(id int) (*entity.User, error) {
	User := entity.User{}

	err := s.db.Where("user_id = $1", id).First(&User).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &User, nil
}

func (s *UserService) Insert(User *entity.User) error {
	err := s.db.Create(User).Error
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "duplicate key value violates unique constraint"):
			return ErrDuplicateRecord
		default:
			return err
		}
	}
	return nil
}

// NOTE: Implement versioning to prevent race update
func (s *UserService) Update(user *entity.User) error {
	// versionbefore := user.Version
	err := s.db.
		Clauses(clause.Returning{}).
		Model(user).
		// Where("version = $11", versionbefore).
		Updates(&entity.User{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Address:   user.Address,
			Phone:     user.Phone,
			Email:     user.Email,
			Password:  user.Password,
			Deposit:   user.Deposit,
			Activated: user.Activated,
			Version:   user.Version + 1},
		).Error
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

func (s *UserService) GetToken(tokenScope, tokenPlaintext string) (*entity.User, error) {
	obj := sha256.Sum256([]byte(tokenPlaintext))
	tokenslice := obj[:]

	var token *entity.Token
	err := s.db.
		Preload("User").
		Where("token_hash = $1 AND scope = $2 AND expiry > $3", tokenslice, tokenScope, time.Now()).
		First(&token).
		Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return token.User, nil
}
