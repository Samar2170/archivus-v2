package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	*gorm.Model
	ID          uint
	Title       string
	Description string
	User        User `gorm:"foreignKey:UserID"`
	UserID      uuid.UUID
}

type Todo struct {
	*gorm.Model
	ID          uint
	Title       string
	Description string
	Status      uint
	Priority    uint
	Project     Project `gorm:"foreignKey:ProjectID"`
	ProjectID   uint
	User        User `gorm:"foreignKey:UserID"`
	UserID      uuid.UUID
}

func (t *Todo) TableName() string {
	return "tempora_todos"
}

func (p *Project) TableName() string {
	return "tempora_projects"
}
