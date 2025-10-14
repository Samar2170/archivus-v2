package syncer

import (
	"archivus-v2/internal"
	"archivus-v2/internal/db"
	"archivus-v2/pkg/logging"
	"errors"
	"time"
)

func runSync() error {
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

	fileErr := runSync()
	errs = append(errs, fileErr)
	return errs
}
