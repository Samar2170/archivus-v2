package synkr

import (
	"archivus-v2/config"
	"archivus-v2/internal/db"
	"archivus-v2/internal/models"
	"archivus-v2/pkg/logging"
	"os"
	"path/filepath"
)

func saveFileMetadata(path string, info os.FileInfo) error {
	var existingFmd models.FileMetadata
	trimmedPath, err := filepath.Rel(config.Config.BaseDir, path)
	if err != nil {
		return err
	}
	err = db.StorageDB.Where("rel_path = ?", trimmedPath).First(&existingFmd).Error
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
}

func syncFilesForDir(dir string) (int64, float64, error) {
	count := int64(0)
	var size float64 = 0
	files, err := os.ReadDir(dir)
	if err != nil {
		return count, size, err
	}
	for _, f := range files {
		info, err := f.Info()
		if err != nil {
			return count, size, err
		}
		size += float64(info.Size())
		if !shouldScanFile(info.Name()) {
			continue
		}
		if count%500 == 0 {
			logging.AuditLogger.Printf("Synced %d files", count)
		}
		if info.IsDir() {
			continue
		} else {
			count++
			err := saveFileMetadata(filepath.Join(dir, f.Name()), info)
			if err != nil {
				return count, size, err
			}
		}
	}
	return count, size, err
}
