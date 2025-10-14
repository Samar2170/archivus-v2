package syncer

import (
	"archivus-v2/internal/db"
	"archivus-v2/internal/models"
	"strings"
)

func shouldScanDir(dir string) bool {
	if dir == "" || dir[0] == '.' || dir[0] == '_' {
		return false
	}
	dirNameSplit := strings.Split(dir, "/")
	lastElement := dirNameSplit[len(dirNameSplit)-1]
	if lastElement[0] == '.' {
		return false
	}
	if _, ok := skipDirsList[lastElement]; ok {
		return false
	}
	return true
}

func createDirEntry(dir *models.Directory) {
	pathSplit := strings.Split(dir.Path, "/")
	dir.Name = pathSplit[len(pathSplit)-1]
	db.StorageDB.Create(&dir)
}
