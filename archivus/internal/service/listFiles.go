package service

import (
	"archivus/config"
	"archivus/internal/db"
	"archivus/internal/models"
	"archivus/internal/utils"
	"archivus/pkg/logging"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type FileMetadataResponse struct {
	ID       uint      `gorm:"primaryKey"`
	Name     string    `gorm:"index"`
	FilePath string    `gorm:"index"`
	UserID   uuid.UUID `gorm:"index"`
	SizeInMb float64
	IsPublic bool

	Tags []models.Tags `gorm:"many2many:file_metadata_tags;"`

	IsImage                    bool
	CompressedVersionAvailable bool `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func FindFiles(userId, search, orderBy, ordering, pageNo string) ([]FileMetadataResponse, error) {
	var files []FileMetadataResponse
	query := db.StorageDB.Model(&models.FileMetadata{}).Where("user_id = ?", userId)
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	if orderBy != "" {
		query = query.Order(fmt.Sprintf("%s %s", orderBy, ordering))
	}
	err := query.Find(&files).Error
	if err != nil {
		return files, utils.HandleError("GetFiles", "Failed to get files for user", err)
	}
	// TODO: Breaking change, pagination is not working
	// results, err := paginateResults(query, pageNo, 50, &files)
	// if err != nil {
	// 	return PaginatedResults{}, utils.HandleError("GetFiles", "Failed to paginate results", err)
	// }
	return files, utils.HandleError("GetFiles", "Failed to get files for user", err)
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

func getRelPaths(folder, username string) (string, string, error) {
	var relPath string
	var folderPath string
	if config.Config.Native {
		var homeDir string
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return relPath, relPath, utils.HandleError("FindFiles", "Failed to get user home directory", err)
		}
		relPath = folder
		folderPath = filepath.Join(homeDir, relPath)
	} else {
		relPath = resolveFolderName(username, folder)
		folderPath = filepath.Join(config.Config.UploadsDir, relPath)
	}
	return relPath, folderPath, nil
}

func GetFiles(userId string, folder string) ([]DirEntry, float64, error) {
	user, err := models.GetUserById(userId)
	if err != nil {
		return nil, 0, utils.HandleError("FindFiles", "Failed to get user by API key", err)
	}
	relPath, folderPath, err := getRelPaths(folder, user.Username)
	if err != nil {
		return nil, 0, utils.HandleError("FindFiles", "Failed to get relative path", err)
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
	fmdMap, err := getFileMetadatas(&files, relPath)
	if err != nil {
		logging.Errorlogger.Error().Msgf("Failed to get file metadatas: %v", err)
	}
	for _, file := range files {
		signedUrl, err := GetSignedUrl(relPath+"/"+file.Name(), user.ID.String())
		if err != nil {
			signedUrl = ""
		}
		fmd, exists := fmdMap[filepath.Join(relPath, file.Name())]
		var fmdId uint
		if exists {
			fmdId = fmd.ID
		} else {
			fmdId = 0
		}
		entries = append(entries, DirEntry{
			ID:        fmdId,
			Name:      file.Name(),
			Path:      relPath + "/" + file.Name(),
			IsDir:     file.IsDir(),
			Extension: filepath.Ext(file.Name()),
			SignedUrl: backendAddr + "/files/download/" + signedUrl,
			Size:      GetSizeForDirEntry(file),
		})
	}
	var folderSize float64
	folderData, err := models.GetDirByPathorName(relPath, folder, user.ID.String())
	if err == nil {
		folderSize = folderData.SizeInMb
	}
	if entries == nil {
		entries = []DirEntry{}
	}
	return entries, folderSize, nil
}
