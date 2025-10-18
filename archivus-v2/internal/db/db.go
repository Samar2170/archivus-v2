package db

import (
	"archivus-v2/config"
	"fmt"
	"log"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var StorageDB *gorm.DB

func getStorageDbFile(dbFile string) string {
	if config.Config.Mode == "prod" {
		if filepath.IsAbs(dbFile) {
			return dbFile
		} else {
			return filepath.Join(config.Config.BaseDir, config.ArchivusV2Dir, dbFile)
		}
	}
	return filepath.Join(config.ProjectBaseDir, dbFile)
}

func InitDB() {
	var err error
	storageDbFile := getStorageDbFile(config.Config.StorageDbFile)
	StorageDB, err = gorm.Open(sqlite.Open(storageDbFile), &gorm.Config{})
	if err != nil {
		fmt.Println(storageDbFile)
		fmt.Println(err)
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
