package models

import (
	"time"

	"gorm.io/gorm"
)

type SyncState struct {
	*gorm.Model
	FilesSynced  int64
	TotalFileMds int64
	TotalDirs    int64

	LastSyncedAt time.Time
}

type DirQueue struct {
	*gorm.Model
	Path         string
	LastSyncedAt time.Time
}
