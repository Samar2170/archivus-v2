package syncer

import (
	"archivus-v2/internal/db"
	"archivus-v2/internal/models"
	"os"
	"path/filepath"
	"time"
)

func ensureQueueHasRoot(root string) error {
	var q models.DirQueue
	err := db.StorageDB.Where("path = ?", root).First(&q).Error
	if err != nil {
		if err.Error() == "record not found" {
			q := models.DirQueue{
				Path:         root,
				LastSyncedAt: time.Now(),
			}
			db.StorageDB.Create(&q)
		}
	}
	return nil
}

func nextDir() (string, error) {
	var dirQueueElement models.DirQueue
	t := time.Now()
	prevYear := t.AddDate(-1, 0, 0)
	err := db.StorageDB.Where("last_synced_at is null or last_synced_at < ?", prevYear).Order("last_synced_at").First(&dirQueueElement).Error
	return dirQueueElement.Path, err
}

func markDirScanned(path string) error {
	tx := db.StorageDB.Begin()
	defer tx.Commit()
	err := tx.Model(&models.DirQueue{}).Where("path = ?", path).Update("last_synced_at", time.Now()).Error
	return err
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

func addSubDirsToQueue(parent string) error {
	entries, err := os.ReadDir(parent)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			if shouldScanDir(entry.Name()) {
				q := models.DirQueue{
					Path: filepath.Join(parent, entry.Name()),
				}
				db.StorageDB.Create(&q)
			}
		}
	}
	return nil
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
