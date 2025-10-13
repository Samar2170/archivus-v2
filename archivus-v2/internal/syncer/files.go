package syncer

import (
	"archivus-v2/config"
	"archivus-v2/internal/db"
	"archivus-v2/internal/models"
	"archivus-v2/pkg/logging"
	"os"
	"path/filepath"
	"strings"
	"time"
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

func syncFilesForDir(dir string) error {
	count := 0
	// TODO: handle files only flat
	// err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
	// 	if err != nil {
	// 		return err
	// 	}
	// 	info, err := d.Info()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if !shouldScanDir(info.Name()) {
	// 		return filepath.SkipDir
	// 	}
	// 	count++
	// 	if count%500 == 0 {
	// 		logging.AuditLogger.Printf("Synced %d files", count)
	// 	}
	// 	if info.IsDir() {
	// 		return nil
	// 	} else {
	// 		return saveFileMetadata(path, info)
	// 	}
	// })
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		info, err := f.Info()
		if err != nil {
			return err
		}
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
				return err
			}
		}
	}
	return err
}

func getSyncState() (time.Time, error) {
	var state models.SyncState
	err := db.StorageDB.Where("id = ?", 1).First(&state).Error
	if err != nil {
		return time.Time{}, err
	}
	return state.LastSyncedAt, nil
}

func setSyncState() error {
	t := time.Now()
	syncState := models.SyncState{
		LastSyncedAt: t,
	}
	err := db.StorageDB.Where("id = ?", 1).FirstOrCreate(&syncState).Error
	if err != nil {
		return err
	}
	return nil
}
func ensureQueueHasRootSubDirs(root string) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			if shouldScanDir(entry.Name()) {
				path := filepath.Join(root, entry.Name())
				q := models.DirQueue{
					Path: path,
				}
				db.StorageDB.Where("path = ?", path).FirstOrCreate(&q)
			}
		}
	}
	return nil
}

func startSync(stop <-chan struct{}) error {
	var err error
	err = ensureQueueHasRoot(config.Config.BaseDir)
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
			for {
				dir, _ := nextDir()
				isShouldScanDir := shouldScanDir(dir)
				if !isShouldScanDir {
					continue
				}
				if err := syncFilesForDir(dir); err != nil {
					return err
				}
				if err := markDirScanned(dir); err != nil {
					return err
				}
				if err := addSubDirsToQueue(dir); err != nil {
					return err
				}
			}

		}
	}
}

var skipDirsList = map[string]bool{
	"venv": true,
	".git": true,
}

func shouldScanDir(dir string) bool {
	if dir == "" || dir[0] == '.' {
		return false
	}
	dirNameSplit := strings.Split(dir, "/")
	lastElement := dirNameSplit[len(dirNameSplit)-1]
	if lastElement[0] == '.' {
		return false
	}
	if _, ok := skipDirsList[lastElement]; ok {
		return false
	}
	return true
}
