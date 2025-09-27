package dirmanager

import (
	"archivus-v2/config"
	"archivus-v2/internal/models"
	"encoding/json"
	"os"
	"path/filepath"
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
