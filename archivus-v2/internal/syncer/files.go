package syncer

import (
	"archivus-v2/config"
	"archivus-v2/internal/db"
	"archivus-v2/internal/models"
	"archivus-v2/pkg/logging"
	"errors"
	"os"
	"path/filepath"
)

var skipDirsList = map[string]bool{
	"venv":         true,
	".git":         true,
	"__pycache__":  true,
	"snap":         true,
	"node_modules": true,
	"bin":          true,
	"build":        true,
}

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

func syncFilesForDir(dir string) (float64, error) {
	count := 0
	var size float64 = 0
	files, err := os.ReadDir(dir)
	if err != nil {
		return size, err
	}
	for _, f := range files {
		info, err := f.Info()
		if err != nil {
			return size, err
		}
		size += float64(info.Size())
		if !shouldScanDir(info.Name()) {
			continue
		}
		count++
		if count%500 == 0 {
			logging.AuditLogger.Printf("Synced %d files", count)
		}
		if info.IsDir() {
			continue
		} else {
			err := saveFileMetadata(filepath.Join(dir, f.Name()), info)
			if err != nil {
				return size, err
			}
		}
	}
	return size, err
}

func formatErrors(errs []error) string {
	var errStr string
	for _, err := range errs {
		errStr += err.Error() + "\n"
	}
	return errStr
}

func startSync(stop <-chan struct{}) error {
	var errs []error
	err := ensureQueueHasRoot(config.Config.BaseDir)
	if err != nil {
		return err
	}
	err = ensureQueueHasRootSubDirs(config.Config.BaseDir)
	if err != nil {
		return err
	}
	for {
		select {
		case <-stop:
			return nil
		default:
			dir, _ := nextDir()
			isShouldScanDir := shouldScanDir(dir)
			if !isShouldScanDir {
				continue
			}
			dirEntry := models.Directory{
				Path: dir,
			}

			size, err := syncFilesForDir(dir)
			if err != nil {
				dirEntry.LastError = err.Error()
				createDirEntry(&dirEntry)
				errs = append(errs, err)
			} else {
				dirEntry.SizeInMb = size / 1024 / 1024
			}
			if err := markDirScanned(dir); err != nil {
				dirEntry.LastError = err.Error()
				createDirEntry(&dirEntry)
				errs = append(errs, err)
			}
			if err := addSubDirsToQueue(dir); err != nil {
				dirEntry.LastError = err.Error()
				createDirEntry(&dirEntry)
				errs = append(errs, err)
			}
			createDirEntry(&dirEntry)
			select {
			case <-stop:
				return errors.New(formatErrors(errs))
			default:
				return errors.New(formatErrors(errs))
			}
		}
	}
}
