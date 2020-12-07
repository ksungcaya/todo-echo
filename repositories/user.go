package repositories

import (
	"github.com/ksungcaya/todo-echo/models"
	"gorm.io/gorm"
)

// UserRepository will interact to the user table.
type UserRepository interface {
	// Methods for querying for single users
	ByID(id uint) *models.User
	ByEmail(email string) *models.User
	ByUsername(username string) *models.User

	// Methods for altering users
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uint) error
}

type userRepoGorm struct {
	db *gorm.DB
}

// little hack to make sure the concrete struct
// correctly implements the UserRepository
var _ UserRepository = &userRepoGorm{}

// NewUserRepository creates instance of UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepoGorm{db}
}

// ByID will look up a user by ID
// If no record was found, the method will return nil
func (ur *userRepoGorm) ByID(id uint) *models.User {
	var u models.User
	err := ur.db.First(&u, id).Error
	if err == nil {
		return &u
	}

	return nil
}

// ByUsername will look up a user by Username
// If no record was found, the method will return nil
func (ur *userRepoGorm) ByUsername(username string) *models.User {
	var u models.User
	err := ur.db.Where(&models.User{Username: username}).First(&u).Error
	if err == nil {
		return &u
	}

	return nil
}

// ByEmail will look up a user by Email
// If no record was found, the method will return nil
func (ur *userRepoGorm) ByEmail(email string) *models.User {
	var u models.User
	err := ur.db.Where(&models.User{Email: email}).First(&u).Error
	if err == nil {
		return &u
	}

	return nil
}

// Create will create a new record to the database
// based on the provided User struct. The password will
// be automatically Hashed here by gorm's BeforeCreate
// hook. For more information, visit models.User.
func (ur *userRepoGorm) Create(user *models.User) error {
	return ur.db.Create(user).Error
}

// Update will update an existing record to the database
// based on the provided User struct. The password "WILL NOT"
// be automatically Hashed here, instead when another struct
// consumes this method, we will check there if a new password
// has been provided and perform the hasing there.
func (ur *userRepoGorm) Update(user *models.User) error {
	return ur.db.Updates(user).Error
}

// Delete will delete a record from the database by ID
func (ur *userRepoGorm) Delete(id uint) error {
	user := models.User{Model: gorm.Model{ID: id}}
	return ur.db.Delete(&user).Error
}
