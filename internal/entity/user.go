package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID    int `gorm:"primaryKey"`
	FirstName string
	LastName  string
	Address   string
	Phone     string
	Email     string
	Password  []byte
	Deposit   float64
	Activated bool
	Version   int `gorm:"default:1"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Orders []*Order `gorm:"foreignKey:UserID;References:UserID"`
}

func (p *User) SetPassword(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.Password = hash
	return nil
}

func (p *User) MatchesPassword(plaintextPassword string) error {
	err := bcrypt.CompareHashAndPassword(p.Password, []byte(plaintextPassword))
	if err != nil {
		return err
	}

	return nil
}
