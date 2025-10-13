package internal

import (
	"archivus-v2/config"
	"archivus-v2/internal/auth"
	"archivus-v2/internal/db"
	"archivus-v2/internal/models"
	"archivus-v2/internal/service/biguploads"
	"os"
)

func Setup(testEnv bool) {
	if !testEnv {
		db.InitDB()
	}
	db.StorageDB.AutoMigrate(
		&models.User{},
		&models.UserPreference{},
		&models.Tags{},
		&models.FileMetadata{},
		&models.Directory{},
	)

	db.StorageDB.AutoMigrate(
		&models.SyncState{},
		&models.DirQueue{},
	)

	checkArchivusDirs()
}

func SetupRun(testEnv bool) {
	Setup(testEnv)
	errors := []error{}
	errors = append(errors, auth.CheckMasterUser(), biguploads.EnsureBigUploadDirs())
	for _, err := range errors {
		if err != nil {
			panic(err)
		}
	}
}

func checkArchivusDirs() {
	if _, err := os.Stat(config.Config.LogsDir); err != nil {
		os.Mkdir(config.Config.LogsDir, os.ModePerm)
	}
}
