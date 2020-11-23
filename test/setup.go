package test

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/ksungcaya/todo-echo/configs"
	"github.com/ksungcaya/todo-echo/database"
	"gorm.io/gorm"
)

// LoadTestEnv will load .env.test
func LoadTestEnv() error {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))

	err := godotenv.Load(fmt.Sprintf("%s/.env.test", filepath.Dir(d)))
	if err != nil {
		log.Fatal("failed to load test env config: ", err)
	}
	return err
}

// InitTestDB will migrate db tables
func InitTestDB() (*gorm.DB, error) {
	if err := LoadTestEnv(); err != nil {
		panic(err)
	}

	config := configs.New()
	db := database.New(&config.Database, false)
	database.AutoMigrate(db)

	if conn, ok := db.DB(); ok != nil {
		defer conn.Close()
	}

	return db, nil
}
