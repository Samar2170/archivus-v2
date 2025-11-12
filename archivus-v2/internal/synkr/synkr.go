package synkr

import (
	"archivus-v2/config"
	"archivus-v2/internal"
	"archivus-v2/internal/models"
	"context"
	"errors"
	"time"
)

func StartSync(ctx context.Context, minutes int) error {
	internal.Setup(false)
	// count files and dirs before syncing
	// count afterwards to get accurate file synced count
	// if no files synced an external source of truth needed to make sure its working and not broken
	initialSyncState, err := getSyncState()
	if err != nil {
		return err
	}
	filesSynced := int64(0)
	stop := make(chan struct{})
	errCh := make(chan error, 1)
	go func() {
		filesSynced, err = sync(stop)
		errCh <- err
	}()

	select {
	case <-ctx.Done():
		_ = setSyncState(filesSynced, initialSyncState.TotalFileMds, initialSyncState.TotalDirs)
		close(stop)
		<-errCh
		return ctx.Err()
	case <-time.After(time.Duration(minutes) * 60 * time.Second):
		_ = setSyncState(filesSynced, initialSyncState.TotalFileMds, initialSyncState.TotalDirs)
		close(stop)
		err := <-errCh
		return err
	}
}

func sync(stop <-chan struct{}) (int64, error) {
	var errs []error
	filesSynced := int64(0)
	err := ensureQueueHasRoot(config.Config.BaseDir)
	if err != nil {
		return filesSynced, err
	}
	for {
		select {
		case <-stop:
			return filesSynced, errors.New(formatErrors(errs))
		default:
			dir, _ := nextDir()
			isShouldScanDir := shouldScanDir(dir)
			if !isShouldScanDir {
				continue
			}
			dirEntry := models.Directory{
				Path: dir,
			}
			count, size, err := syncFilesForDir(dir)
			filesSynced += count
			if err != nil {
				dirEntry.LastError = err.Error()
				errs = append(errs, err)
			} else {
				dirEntry.SizeInMb = size / 1024 / 1024
			}
			err = markDirScanned(dir)
			if err != nil {
				dirEntry.LastError = dirEntry.LastError + ", " + err.Error()
				errs = append(errs, err)
			}
			err = addSubDirsToQueue(dir)
			if err != nil {
				dirEntry.LastError = dirEntry.LastError + ", " + err.Error()
				errs = append(errs, err)
			}
			createDirEntry(&dirEntry)

			select {
			case <-stop:
				return filesSynced, errors.New(formatErrors(errs))
			default:
				continue
			}
		}
	}
}
