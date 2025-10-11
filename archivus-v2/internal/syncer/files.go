package syncer

import (
	"archivus-v2/config"
	"archivus-v2/internal/db"
	"archivus-v2/internal/models"
	"os"
	"path/filepath"
)

func syncFiles() {
	filepath.Walk(config.Config.BaseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		var existingFmd models.FileMetadata
		trimmedPath, err := filepath.Rel(config.Config.BaseDir, path)
		if err != nil {
			return err
		}
		err = db.StorageDB.Where("file_path = ?", trimmedPath).First(&existingFmd).Error
		if err != nil {
			if err.Error() == "record not found" {
				fmd := models.FileMetadata{
					Name:     info.Name(),
					RelPath:  trimmedPath,
					SizeInMb: float64(info.Size()) / 1024 / 1024,
				}
				db.StorageDB.Create(&fmd)
			}
		}
		return nil
	})
}
