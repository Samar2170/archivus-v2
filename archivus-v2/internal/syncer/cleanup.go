package syncer

import (
	"archivus-v2/internal"
	"archivus-v2/internal/db"
	"archivus-v2/internal/models"
	"time"
)

func CleanupDirQueue() {
	internal.Setup(false)
	var dirs []models.DirQueue
	// db.StorageDB.Unscoped().Where("LOWER(path) LIKE ?", "%.idea%").Find(&dirs)
	// var dirs2 []models.DirQueue
	// db.StorageDB.Unscoped().Where("LOWER(path) LIKE ?", "%.vscode%").Find(&dirs2)
	// db.StorageDB.Find(&dirs)

	// for _, dir := range dirs {
	// 	pathSplit := strings.Split(dir.Path, "/")
	// 	lastElement := pathSplit[len(pathSplit)-1]
	// 	if lastElement[0] == '.' || lastElement[0] == '_' || lastElement == "venv" || lastElement == "__pycache__" || lastElement == "node_modules" {
	// 		db.StorageDB.Unscoped().Delete(&dir)
	// 	}
	// 	if len(pathSplit) > 1 {
	// 		lastElement2 := pathSplit[len(pathSplit)-2]
	// 		if lastElement2[0] == '.' || lastElement2[0] == '_' || lastElement2 == "venv" || lastElement2 == "__pycache__" || lastElement2 == "node_modules" {
	// 			db.StorageDB.Unscoped().Delete(&dir)
	// 		}
	// 	}
	// 	if len(pathSplit) > 2 {
	// 		lastElement3 := pathSplit[len(pathSplit)-3]
	// 		fmt.Println(lastElement3, lastElement)
	// 		if lastElement3 == "" {
	// 			continue
	// 		}
	// 		if lastElement3[0] == '.' || lastElement3[0] == '_' || lastElement3 == "venv" || lastElement3 == "__pycache__" || lastElement3 == "node_modules" {
	// 			db.StorageDB.Unscoped().Delete(&dir)
	// 		}
	// 	}
	// }

	var fmds []models.FileMetadata
	// // db.StorageDB.Unscoped().Where("LOWER(rel_path) LIKE ? OR LOWER(rel_path) LIKE ?", "%.android%", "%.cache%").Find(&fmds)
	// // for _, fmd := range fmds {
	// // 	db.StorageDB.Unscoped().Delete(&fmd)
	// // }
	dayBefore := time.Now().Add(-48 * time.Hour)
	db.StorageDB.Unscoped().Where("created_at > ?", dayBefore).Delete(&fmds)
	db.StorageDB.Unscoped().Where("created_at > ?", dayBefore).Delete(&dirs)

	var directories []models.Directory
	db.StorageDB.Unscoped().Where("created_at > ?", dayBefore).Delete(&directories)
}
