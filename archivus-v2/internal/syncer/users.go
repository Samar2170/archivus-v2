package syncer

import (
	"archivus-v2/config"
	"archivus-v2/internal/auth"
	"archivus-v2/internal/models"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

func getCredsData() ([]auth.Credential, error) {
	credsFilePath := filepath.Join(config.Config.BaseDir, config.CredsFile)
	_, err := os.Stat(credsFilePath)
	if os.IsNotExist(err) {
		empty := []auth.Credential{}
		data, _ := json.MarshalIndent(empty, "", " ")
		if err := os.WriteFile(credsFilePath, data, 0644); err != nil {
			return nil, fmt.Errorf("failed to create file: %w", err)
		}
		return empty, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}
	data, err := os.ReadFile(credsFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	var creds []auth.Credential
	err = json.Unmarshal(data, &creds)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal file: %w", err)
	}
	return creds, nil
}

func writeCredsFile(creds []auth.Credential) error {
	// check if creds.json exists
	existingCreds, err := getCredsData()
	if err != nil {
		return err
	}

	// update creds.json
	creds = append(existingCreds, creds...)
	data, err := json.MarshalIndent(creds, "", " ")
	if err != nil {
		return err
	}
	credsFilePath := filepath.Join(config.Config.BaseDir, config.CredsFile)
	return os.WriteFile(credsFilePath, data, 0644)
}

func SyncUsers() []error {
	var errs []error
	var creds []auth.Credential
	dirs, err := os.ReadDir(config.Config.BaseDir)
	if err != nil {
		return []error{err}
	}
	for _, d := range dirs {
		userInfoFile := filepath.Join(config.Config.BaseDir, d.Name(), config.UserInfoFileName)
		_, err := os.Stat(userInfoFile)
		if err == nil {
			c, err := getOrCreateUser(d.Name())
			if err != nil {
				errs = append(errs, err)
				continue
			}
			creds = append(creds, c)
		}
	}
	err = writeCredsFile(creds)
	if err != nil {
		errs = append(errs, err)
	}
	return errs
}

func getOrCreateUser(username string) (auth.Credential, error) {
	_, err := models.GetUserByUsername(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			var c auth.Credential
			c, err = auth.CreateUserByMasterPerm(username)
			if err != nil {
				return auth.Credential{}, err
			}
			return c, nil
		}
	}
	return auth.Credential{}, err
}
