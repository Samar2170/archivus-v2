package syncer

import (
	"archivus/internal/db"
	"archivus/internal/helpers"
	"archivus/internal/service"
	"archivus/pkg/logging"
	"fmt"
)

func Sync() error {
	if db.StorageDB == nil {
		logging.Errorlogger.Error().Msg("StorageDB not initialized")
		fmt.Printf("StorageDB not initialized\n")
		return nil
	}
	err := syncUsers()
	if err != nil {
		logging.Errorlogger.Error().Msgf("Error syncing users: %v", err)
		fmt.Printf("Error syncing users: %s\n", err)
		return err
	}
	err = syncFiles()
	if err != nil {
		logging.Errorlogger.Error().Msgf("Error syncing files: %v", err)
		fmt.Printf("Error syncing files: %s\n", err)
		return err
	}
	return nil
}

func ManualSync() {
	service.Setup(false)
	Sync()
	helpers.UpdateDirsData()
	helpers.UpdateUserDirsData()

}
