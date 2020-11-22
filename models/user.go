package models

import "gorm.io/gorm"

// User model definition
type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(30);unique_index;not null"`
	Email    string `gorm:"type:varchar(100);unique_index;not null"`
	Name     string `gorm:"type:varchar(100);not null"`
	Password string `gorm:"type:varchar(100);"`
}
