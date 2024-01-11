package entity

type Contact struct {
	ContactID int `gorm:"primaryKey"`
	UserID    int
	FirstName string
	LastName  string
	Email     string 
	Address   string
	Phone     string
}
