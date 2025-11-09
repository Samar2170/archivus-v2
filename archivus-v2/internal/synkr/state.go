package synkr

import (
	"archivus-v2/internal/db"
	"archivus-v2/internal/models"
	"time"

	"gorm.io/gorm"
)

func getSyncState() (models.SyncState, error) {
	var state models.SyncState
	err := db.StorageDB.Where("id = ?", 1).First(&state).Error
	if err == gorm.ErrRecordNotFound {
		db.StorageDB.Create(&state)
		err = nil
	}
	return state, err
}

func setSyncState(filesSynced int64, totalFileMds int64, totalDirs int64) error {
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
	if totalFileMds == fmdsCount && totalDirs == dirsCount {
		state.LastErr = "Sync Integrity Check: FAIL"
	}
	err = db.StorageDB.Save(&state).Error

	return err
}
