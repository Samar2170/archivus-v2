package syncer

import (
	"archivus-v2/internal"
	"archivus-v2/internal/db"
	"archivus-v2/internal/models"
	"archivus-v2/pkg/logging"
	"errors"
	"time"
)

func CleanupDirQueue() {
	internal.Setup(false)
	var dirs []models.DirQueue
	// db.StorageDB.Unscoped().Where("LOWER(path) LIKE ?", "%.android%").Find(&dirs)
	// for _, dir := range dirs {
	// 	db.StorageDB.Unscoped().Delete(&dir)
	// }

	var fmds []models.FileMetadata
	// db.StorageDB.Unscoped().Where("LOWER(rel_path) LIKE ? OR LOWER(rel_path) LIKE ?", "%.android%", "%.cache%").Find(&fmds)
	// for _, fmd := range fmds {
	// 	db.StorageDB.Unscoped().Delete(&fmd)
	// }
	dayBefore := time.Now().Add(-24 * time.Hour)
	db.StorageDB.Unscoped().Where("created_at > ?", dayBefore).Delete(&fmds)
	db.StorageDB.Unscoped().Where("created_at > ?", dayBefore).Delete(&dirs)
}

func Sync() []error {
	internal.Setup(false)
	var errs []error
	if db.StorageDB == nil {
		logging.Errorlogger.Error().Msg("Storage DB is nil")
		errs = append(errs, errors.New("storage DB is nil"))
		return errs
	}
	userErrs := SyncUsers()
	errs = append(errs, userErrs...)

	fileErr := startDirSync()
	errs = append(errs, fileErr)
	return errs
}

func startDirSync() error {
	var err error
	stop := make(chan struct{})
	go func() {
		err = startSync(stop)
	}()
	time.Sleep(30 * time.Second)
	close(stop)

	err = setSyncState()
	return err
}
