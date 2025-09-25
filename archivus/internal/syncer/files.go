package syncer

import (
	"archivus/config"
	"archivus/internal/db"
	"archivus/internal/models"
	"os"
	"path/filepath"
)

func syncFileForUser(user models.User) {
	userDir := filepath.Join(config.Config.UploadsDir, user.Username)
	filepath.Walk(userDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		var existingFmd models.FileMetadata
		trimmedPath, err := filepath.Rel(config.Config.UploadsDir, path)
		if err != nil {
			return err
		}
		err = db.StorageDB.Where("user_id = ? AND file_path = ?", user.ID, trimmedPath).First(&existingFmd).Error
		if err != nil {
			if err.Error() == "record not found" {
				fmd := models.FileMetadata{
					Name:     info.Name(),
					FilePath: trimmedPath,
					UserID:   user.ID,
					SizeInMb: float64(info.Size()) / 1024 / 1024,
				}
				db.StorageDB.Create(&fmd)
			}
		}
		return nil
	})

}

func syncFilesUserLevel() error {
	var users []models.User
	err := db.StorageDB.Find(&users).Error
	if err != nil {
		return err
	}
	for _, user := range users {
		syncFileForUser(user)

	}
	return nil

}
