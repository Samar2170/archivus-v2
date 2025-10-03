package service

import (
	"archivus-v2/config"
	"archivus-v2/internal/db"
	"archivus-v2/internal/models"
	"archivus-v2/internal/utils"
	"archivus-v2/pkg/logging"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func GetAllFolders(username string) ([]FolderEntry, error) {
	user, err := models.GetUserByUsername(username)
	if err != nil {
		return nil, utils.HandleError("GetAllFolders", "Failed to get user by API key", err)
	}
	pathFromUploadsDir := filepath.Join(user.Username)
	folderPath := filepath.Join(config.Config.BaseDir, pathFromUploadsDir)
	var subDirs []FolderEntry
	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return utils.HandleError("GetAllFolders", "Failed to walk directory", err)
		}
		if info.IsDir() && path != folderPath {
			subDirs = append(subDirs, FolderEntry{
				Name: info.Name(),
				Path: splitPathTillBaseDir(path),
			})
		}
		return nil
	})
	if err != nil {
		return nil, utils.HandleError("GetAllFolders", "Failed to walk directory", err)
	}

	return subDirs, nil
}

type DirEntry struct {
	ID        uint
	Name      string
	IsDir     bool
	Extension string
	SignedUrl string
	Size      float64
	Path      string

	NavigationPath string
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
		fmdMap[fmd.RelPath] = fmd
	}
	return fmdMap, err
}

func GetFiles(userId string, folder string) ([]DirEntry, float64, error) {
	user, err := models.GetUserById(userId)
	if err != nil {
		return nil, 0, utils.HandleError("FindFiles", "Failed to get user by API key", err)
	}

	// pathFromUploadsDir := resolveFolderName(user.Username, folder)
	var folderPath string
	if user.UserDirLock {
		folderPath = filepath.Join(config.Config.BaseDir, user.Username, folder)
	} else {
		folderPath = filepath.Join(config.Config.BaseDir, folder)
	}
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, 0, utils.HandleError("FindFiles", "Failed to read directory", err)
	}

	var entries []DirEntry
	var backendAddr string

	if config.Config.Mode == "net" {
		currentIp, err := utils.GetPrivateIp()
		if err != nil {
			return nil, 0, utils.HandleError("FindFiles", "Failed to get current IP address", err)
		}
		backendAddr = fmt.Sprintf("%s://%s:%s", config.GetBackendScheme(), currentIp, config.Config.BackendConfig.Port)
	} else {
		backendAddr = fmt.Sprintf("%s://%s", config.GetBackendScheme(), config.GetBackendAddr())
	}
	fmdMap, err := getFileMetadatas(&files, folder)
	if err != nil {
		logging.Errorlogger.Error().Msgf("Failed to get file metadatas: %v", err)
	}
	for _, file := range files {
		signedUrl, err := GetSignedUrl(filepath.Join(folder, file.Name()), user.ID.String())
		if err != nil {
			signedUrl = ""
		}
		fmd, exists := fmdMap[filepath.Join(folder, file.Name())]
		var fmdId uint
		if exists {
			fmdId = fmd.ID
		} else {
			fmdId = 0
		}
		if file.Name()[0] == '.' {
			continue
		}
		entries = append(entries, DirEntry{
			ID:        fmdId,
			Name:      file.Name(),
			Path:      filepath.Join(folder, file.Name()),
			IsDir:     file.IsDir(),
			Extension: filepath.Ext(file.Name()),
			SignedUrl: backendAddr + "/files/download/" + signedUrl,
			Size:      GetSizeForDirEntry(file),
		})
	}
	var folderSize float64
	folderData, err := models.GetDirByPathorName(folder, user.ID.String())
	if err == nil {
		folderSize = folderData.SizeInMb
	}
	if entries == nil {
		entries = []DirEntry{}
	}
	return entries, folderSize, nil
}

func GetSizeForDirEntry(file fs.DirEntry) float64 {
	fi, err := file.Info()
	if err != nil {
		return 0
	}
	return float64(fi.Size() / 1024 / 1024)
}
