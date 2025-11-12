package syncer

import (
	"archivus-v2/internal"
	"archivus-v2/internal/db"
	"archivus-v2/pkg/logging"
	"context"
	"errors"
	"time"
)

func runSync(ctx context.Context) error {
	stop := make(chan struct{})
	errCh := make(chan error, 1)
	// go func() {
	errCh <- startSync(stop)
	// }()
	select {
	case <-ctx.Done():
		close(stop)
		<-errCh
		return ctx.Err()
	case <-time.After(300 * time.Second):
		close(stop)
		err := <-errCh
		_ = setSyncState()
		return err
	}
	// time.Sleep(300 * time.Second)
	// close(stop)

	// err = setSyncState()
	// return err
}

func Sync(ctx context.Context) []error {
	internal.Setup(false)
	var errs []error
	if db.StorageDB == nil {
		logging.Errorlogger.Error().Msg("Storage DB is nil")
		errs = append(errs, errors.New("storage DB is nil"))
		return errs
	}
	// userErrs := SyncUsers()
	// errs = append(errs, userErrs...)

	fileErr := runSync(ctx)
	errs = append(errs, fileErr)
	return errs
}
