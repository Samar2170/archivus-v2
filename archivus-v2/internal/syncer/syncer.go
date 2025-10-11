package syncer

import (
	"archivus-v2/internal/db"
	"archivus-v2/pkg/logging"
)

func Sync() []error {
	var errs []error
	if db.StorageDB == nil {
		logging.Errorlogger.Error().Msg("Storage DB is nil")
		return nil
	}
	userErrs := SyncUsers()
	errs = append(errs, userErrs...)
	return errs
}
