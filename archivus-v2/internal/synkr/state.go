package synkr

import (
	"archivus-v2/internal/db"
	"archivus-v2/internal/models"
	"time"
)

func getSyncState() (models.SyncState, error) {
	var state models.SyncState
	err := db.StorageDB.Where("id = ?", 1).First(&state).Error
	return state, err
}

func setSyncState(filesSynced int64) error {
	var fmdsCount int64
	var dirsCount int64
	db.StorageDB.Model(&models.FileMetadata{}).Count(&fmdsCount)
	db.StorageDB.Model(&models.Directory{}).Count(&dirsCount)

	var state models.SyncState
	err := db.StorageDB.Where("id = ?", 1).First(&state).Error
	if err != nil {
		state = models.SyncState{
			LastSyncedAt: time.Now(),
		}
		db.StorageDB.Create(&state)
	}
	state.FilesSynced = filesSynced
	state.TotalFileMds = fmdsCount
	state.TotalDirs = dirsCount
	state.LastSyncedAt = time.Now()
	err = db.StorageDB.Save(&state).Error

	return err
}
