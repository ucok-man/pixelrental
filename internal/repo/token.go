package repo

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"errors"
	"time"

	"github.com/ucok-man/pixelrental/internal/entity"
	"gorm.io/gorm"
)

type TokenService struct {
	db *gorm.DB
}

func GenerateToken(userID int, ttl time.Duration, scope string) (*entity.Token, error) {
	token := &entity.Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.TokenHash = hash[:]

	return token, nil
}

func (s TokenService) GenerateAndInsert(userID int, ttl time.Duration, scope string) (*entity.Token, error) {
	token, err := GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = s.Insert(token)
	return token, err
}

func (s TokenService) Insert(token *entity.Token) error {
	err := s.db.Create(token).Error
	if err != nil {
		return err
	}
	return nil
}

func (s TokenService) GetTokenInvoice(scope string, userID int) (*entity.Token, error) {
	var token *entity.Token
	err := s.db.Where("scope = $1 AND user_id = $2 AND expiry > $3", scope, userID, time.Now()).First(&token).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return token, nil
}

func (s TokenService) DeleteTokenAll(scope string, userID int) error {
	err := s.db.Where("scope = $1 AND user_id = $2", scope, userID).Delete(&entity.Token{}).Error
	if err != nil {
		return err
	}
	return nil
}
