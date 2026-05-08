package repositories

import (
	"echo_practice/internal/models"
	"errors"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var u models.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var u models.User
	if err := r.db.Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var u models.User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
