package service

import (
	"archivus-v2/config"
	"archivus-v2/internal/db"
	"archivus-v2/internal/models"
	"archivus-v2/internal/utils"
	"archivus-v2/pkg/logging"
	"io"
	"os"
	"path/filepath"
)

func MoveFile(userId, filePath, destFolder string, copy bool) error {
	// Implement the logic to move the file to the specified folder
	// This is a placeholder for the actual implementation
	var fileMd models.FileMetadata
	err := db.StorageDB.Model(&models.FileMetadata{}).Where("rel_path = ? AND uploaded_by_id = ?", filePath, userId).First(&fileMd).Error
	if err != nil {
		logging.Errorlogger.Error().Msgf("Error retrieving file metadata: %v", err)
		return utils.HandleError("Move File", "Failed to retrieve file metadata", err)
	}
	tx := db.StorageDB.Begin()
	newFilePath := filepath.Join(config.Config.BaseDir, destFolder, fileMd.Name)
	newFile, err := os.Create(newFilePath)
	if err != nil {
		tx.Rollback()
		logging.Errorlogger.Error().Msgf("Error creating new file: %v", err)
		return utils.HandleError("Move File", "Failed to create new file", err)
	}
	defer newFile.Close()
	oldFile, err := os.Open(filepath.Join(config.Config.BaseDir, fileMd.RelPath))
	if err != nil {
		tx.Rollback()
		logging.Errorlogger.Error().Msgf("Error opening old file: %v", err)
		return utils.HandleError("Move File", "Failed to open old file", err)
	}
	_, err = io.Copy(newFile, oldFile)
	if err != nil {
		tx.Rollback()
		logging.Errorlogger.Error().Msgf("Error moving file: %v", err)
		return utils.HandleError("Move File", "Failed to move file", err)
	}
	if !copy {
		err = os.Remove(filepath.Join(config.Config.BaseDir, fileMd.RelPath))
		if err != nil {
			tx.Rollback()
			logging.Errorlogger.Error().Msgf("Error removing old file: %v", err)
			return utils.HandleError("Move File", "Failed to remove old file", err)
		}
	}
	fileMd.RelPath = newFilePath
	tx.Save(&fileMd)
	if err := tx.Commit().Error; err != nil {
		logging.Errorlogger.Error().Msgf("Error committing transaction: %v", err)
		return utils.HandleError("Move File", "Failed to commit transaction", err)
	}
	return nil
}

func DeleteFile(userId, fileId string) error {
	var fileMd models.FileMetadata
	err := db.StorageDB.Model(&models.FileMetadata{}).Where("rel_path = ? AND uploaded_by_id = ?", fileId, userId).First(&fileMd).Error
	if err != nil {
		logging.Errorlogger.Error().Msgf("Error retrieving file metadata: %v", err)
		return utils.HandleError("Delete File", "Failed to retrieve file metadata", err)
	}

	tx := db.StorageDB.Begin()
	if err := tx.Delete(&fileMd).Error; err != nil {
		tx.Rollback()
		logging.Errorlogger.Error().Msgf("Error deleting file metadata: %v", err)
		return utils.HandleError("Delete File", "Failed to delete file metadata", err)
	}

	if err := os.Remove(filepath.Join(config.Config.BaseDir, fileMd.RelPath)); err != nil {
		tx.Rollback()
		logging.Errorlogger.Error().Msgf("Error removing file from storage: %v", err)
		return utils.HandleError("Delete File", "Failed to remove file from storage", err)
	}

	if err := tx.Commit().Error; err != nil {
		logging.Errorlogger.Error().Msgf("Error committing transaction: %v", err)
		return utils.HandleError("Delete File", "Failed to commit transaction", err)
	}

	return nil
}
