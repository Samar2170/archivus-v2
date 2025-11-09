package synkr

import (
	"archivus-v2/internal/db"
	"archivus-v2/internal/models"
	"strings"
)

func createDirEntry(dir *models.Directory) {
	pathSplit := strings.Split(dir.Path, "/")
	dir.Name = pathSplit[len(pathSplit)-1]
	db.StorageDB.Create(&dir)
}
