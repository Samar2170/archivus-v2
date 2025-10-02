package server

import (
	"archivus-v2/internal/models"
	"archivus-v2/pkg/logging"
	"strings"
)

func checkUserWriteAccess(username string) (bool, error) {
	user, err := models.GetUserByUsername(username)
	if err != nil {
		logging.Errorlogger.Error().Msg(err.Error())
		return false, err
	}
	if !user.WriteAccess {
		return false, nil
	}
	return true, nil
}

func checkUserWriteAccessToFolder(username string, folder string) (bool, error) {
	user, err := models.GetUserByUsername(username)
	if err != nil {
		logging.Errorlogger.Error().Msg(err.Error())
		return false, err
	}
	if !user.WriteAccess {
		return false, nil
	}
	folderSplit := strings.Split(folder, "/")
	if user.UserDirLock {
		if folderSplit[0] == username || folderSplit[0] == "" {
			return true, nil
		}
		return false, nil
	}
	return true, nil
}
