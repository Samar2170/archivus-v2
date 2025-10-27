package service

import (
	"archivus-v2/config"
	"archivus-v2/internal/db"
	"archivus-v2/internal/models"
	"archivus-v2/internal/utils"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func GetSignedUrl(filePath string, userId string) (string, error) {
	var fmd models.FileMetadata
	db.StorageDB.Where("rel_path = ? AND uploaded_by_id = ?", filePath, userId).First(&fmd)

	expiresAt := time.Now().Add(600 * time.Minute).Unix()
	expriresAtStr := fmt.Sprintf("%d", expiresAt)
	content := fmt.Sprintf("%s%s%s", config.Config.SecretKey, expriresAtStr, filePath)

	signature := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%s?signature=%s&expires_at=%d", filePath, hex.EncodeToString(signature[:]), expiresAt), nil
}

func DownloadFile(filePath, signature, expiresAt string, compressed bool) ([]byte, error) {
	var absPath string
	content := fmt.Sprintf("%s%s%s", config.Config.SecretKey, expiresAt, filePath)
	hash := sha256.Sum256([]byte(content))
	if hex.EncodeToString(hash[:]) != signature {
		return nil, utils.HandleError("DownloadFile", "Invalid signature", errors.New("invalid signature"))
	}
	absPath = filepath.Join(config.Config.BaseDir, filePath)
	f, err := os.ReadFile(absPath)
	if err != nil {
		return nil, utils.HandleError("DownloadFile", "Failed to read file", err)
	}
	return f, nil
}
