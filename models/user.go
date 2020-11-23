package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User model definition
type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(30);unique_index;not null"`
	Email    string `gorm:"type:varchar(100);unique_index;not null"`
	Name     string `gorm:"type:varchar(100);not null"`
	Password string `gorm:"type:varchar(100);"`
}

// BeforeCreate is called before saving the new user to the database
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if hashed, err := HashPassword(u.Password); err == nil {
		tx.Statement.SetColumn("Password", hashed)
	}

	return
}

// HashPassword hashes password using Bcrypt
func HashPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", errors.New("Password should not be empty")
	}

	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(h), err
}

// CheckPassword compares plain text password to the
// current User type's password u.Password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
