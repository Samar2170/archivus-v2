package service

import (
	"archivus-v2/internal/db"
	"archivus-v2/internal/models"
	"archivus-v2/internal/utils"
	"bufio"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

func SaveFile(file multipart.File, fileHeader *multipart.FileHeader, username, folderPath, tags string) error {
	user, err := models.GetUserByUsername(username)
	if err != nil {
		return utils.HandleError("SaveFile", "Failed to get user by username", err)
	}
	folderPathSplit := strings.Split(folderPath, "/")
	if user.UserDirLock && folderPathSplit[0] != user.Username {
		return utils.HandleError("SaveFile", "Invalid folder path", nil)
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
	tx := db.StorageDB.Begin()

	tagsArray, err := GetTags(tags, tx)
	if err != nil {
		tx.Rollback()
		return utils.HandleError("SaveFile", "Failed to get tags", err)
	}
	fmd := models.FileMetadata{
		Name:         fileHeader.Filename,
		AbsPath:      pathFromUploadsDir,
		UploadedByID: user.ID,
		SizeInMb:     float64(fileHeader.Size) / 1024 / 1024,
	}
	fmd.Tags = tagsArray
	err = tx.Create(&fmd).Error
	if err != nil {
		tx.Rollback()
		return utils.HandleError("SaveFile", "Failed to create file metadata record", err)
	}
	reader := bufio.NewReader(file)
	writer := io.Writer(f)
	_, err = io.Copy(writer, reader)
	if err != nil {
		tx.Rollback()
		return utils.HandleError("SaveFile", "Failed to copy file content", err)
	}
	tx.Commit()
	return nil
}
