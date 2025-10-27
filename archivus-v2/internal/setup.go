package internal

import (
	"archivus-v2/config"
	"archivus-v2/internal/auth"
	"archivus-v2/internal/db"
	"archivus-v2/internal/models"
	"archivus-v2/internal/service/biguploads"
	"os"
	"path/filepath"
)

func Setup(testEnv bool) {
	checkArchivusDirs()

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

	db.StorageDB.AutoMigrate(
		&models.Project{},
		&models.Todo{},
	)

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
	if config.Config.Mode == "prod" {
		if _, err := os.Stat(filepath.Join(config.Config.BaseDir, config.ArchivusV2Dir)); err != nil {
			os.Mkdir(filepath.Join(config.Config.BaseDir, config.ArchivusV2Dir), os.ModePerm)
		}
	}
	if _, err := os.Stat(config.Config.LogsDir); err != nil {
		os.Mkdir(config.Config.LogsDir, os.ModePerm)
	}

}
