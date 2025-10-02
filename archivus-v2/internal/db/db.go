package db

import (
	"archivus-v2/config"
	"log"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var StorageDB *gorm.DB

func InitDB() {
	var err error
	storageDbFile := filepath.Join(config.ProjectBaseDir, config.Config.StorageDbFile)
	StorageDB, err = gorm.Open(sqlite.Open(filepath.Join(config.ProjectBaseDir, storageDbFile)), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
}

func InitTestDB() {
	var err error
	testDbFilePath := filepath.Join(config.ProjectBaseDir, config.TestDbFile)
	StorageDB, err = gorm.Open(sqlite.Open(testDbFilePath), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect test database")
	}
}
