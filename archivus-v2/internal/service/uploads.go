package service

import (
	"archivus-v2/internal/db"
	"archivus-v2/internal/models"
	"archivus-v2/internal/utils"
	"archivus-v2/pkg/logging"
	"bufio"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func SaveFileMetaData(tags, filename, pathFromUploadsDir string, userID uuid.UUID, fileSize float64) error {
	tx := db.StorageDB.Begin()
	tagsArray, err := GetTags(tags, tx)
	if err != nil {
		tx.Rollback()
		return utils.HandleError("SaveFile", "Failed to get tags", err)
	}
	fmd := models.FileMetadata{
		Name:         filename,
		RelPath:      pathFromUploadsDir,
		AbsPath:      pathFromUploadsDir,
		UploadedByID: userID,
		SizeInMb:     float64(fileSize) / 1024 / 1024,
	}
	fmd.Tags = tagsArray
	err = tx.Create(&fmd).Error
	if err != nil {
		tx.Rollback()
		return utils.HandleError("SaveFile", "Failed to create file metadata record", err)
	}
	tx.Commit()
	return nil
}

func SaveFile(file multipart.File, fileHeader *multipart.FileHeader, username, folderPath, tags string) error {
	user, err := models.GetUserByUsername(username)
	if err != nil {
		return utils.HandleError("SaveFile", "Failed to get user by username", err)
	}
	if user.UserDirLock {
		folderPath = filepath.Join(user.Username, folderPath)
	}
	filePath, pathFromUploadsDir, err := getNewFileName(fileHeader.Filename, folderPath)
	if err != nil {
		return utils.HandleError("SaveFile", "Failed to get new file name", err)
	}
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return utils.HandleError("SaveFile", "Failed to open file for writing", err)
	}
	defer f.Close()
	reader := bufio.NewReader(file)
	writer := io.Writer(f)
	_, err = io.Copy(writer, reader)
	if err != nil {
		return utils.HandleError("SaveFile", "Failed to copy file content", err)
	}
	err = SaveFileMetaData(tags, fileHeader.Filename, pathFromUploadsDir, user.ID, float64(fileHeader.Size))
	if err != nil {
		logging.Errorlogger.Println("Failed to save file metadata:", err)
	}
	return nil
}
