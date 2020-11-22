package repositories

import (
	"github.com/ksungcaya/todo-echo/models"
	"gorm.io/gorm"
)

// UserRepository will interact to the user table.
type UserRepository interface {
	// Methods for querying for single users
	ByID(id uint) (*models.User, error)
	ByEmail(email string) (*models.User, error)
	ByUsername(username string) (*models.User, error)

	// Methods for altering users
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uint) error
}

type userRepoGorm struct {
	db *gorm.DB
}

var _ UserRepository = &userRepoGorm{}

// NewUserRepository creates instance of UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepoGorm{db}
}

func (ur *userRepoGorm) ByID(id uint) (*models.User, error) {
	return &models.User{}, nil
}

func (ur *userRepoGorm) ByEmail(email string) (*models.User, error) {
	return &models.User{}, nil
}

func (ur *userRepoGorm) ByUsername(username string) (*models.User, error) {
	return &models.User{}, nil
}

func (ur *userRepoGorm) Create(user *models.User) error {
	return ur.db.Create(user).Error
}

func (ur *userRepoGorm) Update(user *models.User) error {
	return nil
}

func (ur *userRepoGorm) Delete(id uint) error {
	return nil
}
