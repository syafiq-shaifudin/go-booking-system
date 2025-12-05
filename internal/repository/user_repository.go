package repository

import (
	"go-booking-system/internal/domain"

	"gorm.io/gorm"
)

// UserRepository defines data access methods for User
type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id uint) (*domain.User, error)
	FindByUUID(uuid string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id uint) error
}

// userRepository implements UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create inserts a new user into database
func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

// FindByEmail retrieves user by email address
func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID retrieves user by primary key ID
func (r *userRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUUID retrieves user by UUID
func (r *userRepository) FindByUUID(uuid string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("uuid = ?", uuid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update saves user changes to database
func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

// Delete soft deletes a user by ID
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&domain.User{}, id).Error
}
