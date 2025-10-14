package models

import (
	"time"

	"gorm.io/gorm"
)

type SyncState struct {
	*gorm.Model
	LastSyncedAt time.Time
}

type DirQueue struct {
	*gorm.Model
	Path         string
	LastSyncedAt time.Time
}
