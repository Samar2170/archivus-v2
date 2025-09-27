package service

import (
	"archivus/config"
	"archivus/internal/db"
	"archivus/internal/models"
	"archivus/pkg/logging"
	"log"
	"os"
)

func Setup(testEnv bool) {
	if !testEnv {
		db.InitDB()
	}
	// This function is intentionally left empty.
	// It can be used to initialize any global variables or configurations
	// that are needed for the service to run.
	// Currently, it serves as a placeholder for future setup logic.
	db.StorageDB.AutoMigrate(
		&models.User{},
		&models.UserPreference{},
		&models.Tags{},
		&models.FileMetadata{},
		&models.Directory{},
	)
	CheckMasterUser()
	CheckArchivusDirs()
}

func CheckArchivusDirs() {
	if _, err := os.Stat(config.Config.UploadsDir); err != nil {
		os.Mkdir(config.Config.UploadsDir, os.ModePerm)
	}
	if _, err := os.Stat(config.Config.LogsDir); err != nil {
		os.Mkdir(config.Config.UploadsDir, os.ModePerm)
	}
}

func CheckMasterUser() {
	var count int64
	db.StorageDB.Model(&models.User{}).Where("is_master = ?", true).Count(&count)
	if count == 0 {
		// No master user found, create one
		logging.Errorlogger.Println("No master user found. Please create a master user using the 'new-user' command.")
		log.Fatal("No master user found. Please create a master user using the 'new-user' command.")
	}
}
