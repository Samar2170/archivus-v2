package syncer

import (
	"archivus/config"
	"archivus/internal/auth"
	"archivus/internal/models"
	"archivus/pkg/logging"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func getCredsData() ([]auth.Credential, error) {
	credsFilePath := filepath.Join(config.BaseDir, config.CredsFile)
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
	credsFilePath := filepath.Join(config.BaseDir, config.CredsFile)
	return os.WriteFile(credsFilePath, data, 0644)
}

func syncUsers() error {
	userFolders, err := os.ReadDir(config.Config.UploadsDir)
	if err != nil {
		return err
	}
	var creds []auth.Credential
	for _, userFolder := range userFolders {
		if userFolder.IsDir() {
			if _, err := models.GetUserByUsername(userFolder.Name()); err != nil {
				c, err := auth.CreateUserByMasterPerm(userFolder.Name())
				if err != nil {
					logging.Errorlogger.Error().Msgf("Error creating user for %s: %v", userFolder.Name(), err)
					continue
				}
				creds = append(creds, c)
			}
		}
	}
	err = writeCredsFile(creds)
	return err
}
