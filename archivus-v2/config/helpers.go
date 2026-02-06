package config

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	// Encode to base64 to ensure it's printable and URL-safe
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func GetBackendBindAddr() string {
	if Config.BackendConfig.BindAddr != "" {
		return Config.BackendConfig.BindAddr + ":" + Config.BackendConfig.Port
	}
	return GetBackendAddr()
}

func GetBackendAddr() string {
	return Config.BackendConfig.BaseUrl + ":" + Config.BackendConfig.Port
}

func GetBackendScheme() string {
	return Config.BackendConfig.Scheme
}
func GetFrontendScheme() string {
	return Config.FrontEndConfig.Scheme
}

func GetFrontendAddr() string {
	return fmt.Sprintf("%s:%s", Config.FrontEndConfig.BaseUrl, Config.FrontEndConfig.Port)
}

func createOrLoadServerSalt() (string, error) {
	saltFilePath := filepath.Join(ProjectBaseDir, ".server_salt.txt")
	if _, err := os.Stat(saltFilePath); os.IsNotExist(err) {
		// File does not exist, create a new salt
		newSalt, err := generateRandomString(16)
		if err != nil {
			return "", err
		}
		err = os.WriteFile(saltFilePath, []byte(newSalt), 0600)
		if err != nil {
			log.Fatal("Failed to create server salt file:", err)
		}
		return newSalt, nil
	}
	// File exists, read the salt
	saltBytes, err := os.ReadFile(saltFilePath)
	if err != nil {
		log.Fatal("Failed to read server salt file:", err)
	}
	return string(saltBytes), nil
}
