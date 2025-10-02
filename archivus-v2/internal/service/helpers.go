package service

import (
	"archivus-v2/config"
	"archivus-v2/internal/models"
	"archivus-v2/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
)

func GetTags(tagString string, tx *gorm.DB) ([]models.Tags, error) {
	var err error
	tagsSplit := strings.Split(tagString, ",")
	var tagsArray []models.Tags
	for _, tag := range tagsSplit {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		var tagModel models.Tags
		// Check if tag already exists
		err = tx.Where("tag = ?", tag).First(&tagModel).Error
		if err != nil {
			if err.Error() == "record not found" {
				tagModel = models.Tags{Tag: tag}
				err = tx.Create(&tagModel).Error
				if err != nil {
					tx.Rollback()
					return tagsArray, utils.HandleError("GetTags", "Failed to create tag", err)
				}
			}
		}
		tagsArray = append(tagsArray, tagModel)
	}
	if len(tagsArray) == 0 {
		tagsArray = nil // No tags provided, set to nil
	}
	return tagsArray, nil
}

func getAbsoluteFilePath(fileName, folderPath string) (string, string, error) {
	filenameCleand := strings.Replace(fileName, " ", "_", -1)
	pathFromBaseDir := filepath.Join(folderPath, filenameCleand)
	filePath := filepath.Join(config.Config.BaseDir, pathFromBaseDir)
	return filePath, pathFromBaseDir, nil
}

func getNewFileName(fileName, folderPath string) (string, string, error) {
	absFilePath, pathFromUploadsDir, err := getAbsoluteFilePath(fileName, folderPath)
	if err != nil {
		return "", pathFromUploadsDir, utils.HandleError("getNewFileName", "Failed to get absolute file path", err)
	}
	if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
		return absFilePath, pathFromUploadsDir, nil
	} else {
		// File already exists, generate a new name
		ext := filepath.Ext(fileName)
		baseName := strings.TrimSuffix(fileName, ext)
		newFileName := fmt.Sprintf("%s_%d%s", baseName, time.Now().Unix(), ext)
		return getAbsoluteFilePath(newFileName, folderPath)
	}
}

type FolderEntry struct {
	Name string
	Path string
}

func splitPathTillBaseDir(path string) string {
	split := strings.Split(path, "/")
	for i := len(split) - 1; i >= 0; i-- {
		if split[i] == config.Config.BaseDir {
			return filepath.Join(split[i+1:]...)
		}
	}
	return path
}
