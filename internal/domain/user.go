package domain

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Email           string         `gorm:"uniqueIndex;not null" json:"email"`
	Password        string         `gorm:"not null" json:"-"` // "-" means don't return in JSON
	Name            string         `gorm:"not null" json:"name"`
	MobileCountryId *uint          `gorm:"default:null"`
	Phone           string         `json:"phone"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	UUID            string         `gorm:"not null" json:"uuid"`
}

// Hash password before saving
func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Check password
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	// if u.UUID == "" {
	// 	u.UUID = uuid.New().String()
	// }
	u.UUID = uuid.New().String()
	return nil
}
