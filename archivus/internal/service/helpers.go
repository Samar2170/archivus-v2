package service

import (
	"archivus/internal/db"
	"archivus/internal/models"
	"archivus/pkg/logging"
	"os"
	"path/filepath"
	"strings"
)

func getFileMetadatas(files *[]os.DirEntry, pathFromUploadsDir string) (map[string]models.FileMetadata, error) {
	paths := []string{}
	fmdMap := make(map[string]models.FileMetadata)
	for _, file := range *files {
		if file.IsDir() {
			// Skip directories
			continue
		}
		paths = append(paths, filepath.Join(pathFromUploadsDir, file.Name()))
	}
	var fmds []models.FileMetadata
	err := db.StorageDB.Model(&models.FileMetadata{}).Where("file_path IN ?", paths).Find(&fmds).Error
	if err != nil {
		logging.Errorlogger.Error().Msgf("Failed to get file metadata for paths: %v, error: %v", paths, err)
		return fmdMap, err
	}
	for _, fmd := range fmds {
		fmdMap[fmd.FilePath] = fmd
	}
	return fmdMap, err
}

func resolveFolderName(username, folder string) string {
	if folder == "" || folder == "/" {
		return username
	}
	folderSplit := strings.Split(folder, "/")
	if folderSplit[0] == username {
		return folder
	}
	return filepath.Join(username, folder)
}
