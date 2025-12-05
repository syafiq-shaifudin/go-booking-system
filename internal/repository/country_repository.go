package repository

import (
	"go-booking-system/internal/domain"

	"gorm.io/gorm"
)

// CountryRepository defines data access methods for Country
type CountryRepository interface {
	FindByShortname(shortname string) (*domain.Country, error)
	FindByID(id uint) (*domain.Country, error)
	FindAll() ([]domain.Country, error)
}

// countryRepository implements CountryRepository
type countryRepository struct {
	db *gorm.DB
}

// NewCountryRepository creates a new country repository instance
func NewCountryRepository(db *gorm.DB) CountryRepository {
	return &countryRepository{db: db}
}

// FindByShortname retrieves country by shortname (e.g., "US", "UK")
func (r *countryRepository) FindByShortname(shortname string) (*domain.Country, error) {
	var country domain.Country
	err := r.db.Where("shortname = ?", shortname).First(&country).Error
	if err != nil {
		return nil, err
	}
	return &country, nil
}

// FindByID retrieves country by primary key ID
func (r *countryRepository) FindByID(id uint) (*domain.Country, error) {
	var country domain.Country
	err := r.db.First(&country, id).Error
	if err != nil {
		return nil, err
	}
	return &country, nil
}

// FindAll retrieves all countries
func (r *countryRepository) FindAll() ([]domain.Country, error) {
	var countries []domain.Country
	err := r.db.Find(&countries).Error
	return countries, err
}
