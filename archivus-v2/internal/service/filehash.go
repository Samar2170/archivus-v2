package service

import (
	"encoding/hex"
	"io"
	"os"

	"archivus-v2/internal/db"
	"archivus-v2/internal/models"

	"github.com/google/uuid"
	"github.com/zeebo/blake3"
)

func HashFileBlake3(path string) (string, int64, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	hasher := blake3.New()
	size, err := io.Copy(hasher, file)
	if err != nil {
		return "", 0, err
	}

	sum := hasher.Sum(nil)
	return hex.EncodeToString(sum), size, nil
}

func CreateFileHash(filePath string, fileMetadataID uint, userID uuid.UUID) error {
	hashStr, size, err := HashFileBlake3(filePath)
	if err != nil {
		return err
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	hashRecord := models.FileContentHash{
		ID:             uuid.New(),
		UserID:         userID,
		FileMetadataID: fileMetadataID,
		Path:           filePath,
		Size:           size,
		ModTime:        fileInfo.ModTime(),
		Hash:           hashStr,
	}

	return db.StorageDB.Create(&hashRecord).Error
}
