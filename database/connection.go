package database

import (
	"fmt"

	"os"

	"github.com/ksungcaya/todo-echo/configs"
	"github.com/ksungcaya/todo-echo/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// New creates new db instance
func New(config *configs.DatabaseConfig, isProd bool) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.MySQL.Username,
		config.MySQL.Password,
		config.MySQL.Host,
		config.MySQL.Port,
		config.MySQL.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	return db
}

// TestDB creates a test database
func TestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("./../todo_test.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("storage err: ", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(3)

	return db
}

// DropTestDB drops test database
func DropTestDB() error {
	if err := os.Remove("./../todo_test.db"); err != nil {
		return err
	}
	return nil
}

//AutoMigrate migrates models
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
	)
}

// Refresh drops all tables and rebuilds them
func Refresh(db *gorm.DB) error {
	err := db.Migrator().DropTable(&models.User{})
	if err != nil {
		return err
	}

	return AutoMigrate(db)
}
