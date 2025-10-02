package dirmanager

import (
	"archivus-v2/config"
	"archivus-v2/internal/models"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func CreateDirForUser(user models.User) error {
	info, err := os.Stat(filepath.Join(config.Config.BaseDir, user.Username))
	if err != nil {
		if os.IsNotExist(err) {
			return os.Mkdir(filepath.Join(config.Config.BaseDir, user.Username), os.ModePerm)
		}
		return err
	}
	if !info.IsDir() {
		return os.Mkdir(filepath.Join(config.Config.BaseDir, user.Username), os.ModePerm)
	}
	userInfoData := map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
		"id":       user.ID,
	}
	userInfoFilePath := filepath.Join(config.Config.BaseDir, user.Username, config.UserInfoFileName)
	f, err := os.Create(userInfoFilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(userInfoData)
	return err
}

func createFolder(dirPath string) error {
	isAbs := filepath.IsAbs(dirPath)
	if !isAbs {
		dirPath = filepath.Join(config.Config.BaseDir, dirPath)
	}
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}
	} else if info.IsDir() {
		return nil
	} else {
		return err
	}
	return nil
}

func CreateFolder(subFolder, username string) error {
	var dirPath string
	if subFolder == "" {
		return errors.New("subFolder cannot be empty")
	}
	splitPath := strings.Split(subFolder, "/")
	user, err := models.GetUserByUsername(username)
	if err != nil {
		return err
	}
	if user.UserDirLock {
		if splitPath[0] == username {
			dirPath = filepath.Join(config.Config.BaseDir, subFolder)
			return createFolder(dirPath)
		}
		dirPath = filepath.Join(config.Config.BaseDir, username, subFolder)
	} else {
		dirPath = filepath.Join(config.Config.BaseDir, subFolder)
	}
	return createFolder(dirPath)
}
