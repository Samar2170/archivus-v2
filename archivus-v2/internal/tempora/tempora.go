package tempora

import (
	"archivus-v2/internal/db"
	"archivus-v2/internal/models"
	"strconv"

	"github.com/google/uuid"
)

const (
	todo = iota
	inProgress
	done
)

const (
	low = iota
	medium
	high
)

type TodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      uint   `json:"status"`
	Priority    uint   `json:"priority"`
	ProjectID   uint   `json:"projectId"`
}

type UpdateTodosRequest struct {
	Id     uint
	Status uint
}
type ListProjectResponseItem struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ProjectId   uint   `json:"projectId"`
	UserId      string `json:"userId"`
}

type TodoResponseItem struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      uint   `json:"status"`
	Priority    uint   `json:"priority"`
	ProjectID   uint   `json:"projectId"`
}

func CreateTodos(todos []TodoRequest, userId string) error {
	var newTodos []models.Todo
	var err error
	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		return err
	}
	for _, todo := range todos {
		newTodo := models.Todo{
			Title:       todo.Title,
			Description: todo.Description,
			Status:      todo.Status,
			Priority:    todo.Priority,
			ProjectID:   todo.ProjectID,
			UserID:      userIdUUID,
		}
		newTodos = append(newTodos, newTodo)
	}
	return db.StorageDB.Create(&newTodos).Error
}

func ListTodos(userId, projectId string) ([]TodoResponseItem, error) {
	var todos []models.Todo
	if projectId == "" {
		db.StorageDB.Where("user_id = ?", userId).Preload("Project").Find(&todos)
	} else {
		projectIdInt, err := strconv.Atoi(projectId)
		if err != nil {
			return nil, err
		}
		db.StorageDB.Where("user_id = ? AND project_id = ?", userId, projectIdInt).Preload("Project").Find(&todos)
	}
	var response []TodoResponseItem
	for _, todo := range todos {
		response = append(response, TodoResponseItem{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
			Priority:    todo.Priority,
			Status:      todo.Status,
			ProjectID:   todo.ProjectID,
		})
	}
	return response, nil
}

func UpdateTodos(updateTodos []UpdateTodosRequest) error {
	tx := db.StorageDB.Begin()
	for _, todo := range updateTodos {
		tx.Model(&models.Todo{}).Where("id = ?", todo.Id).Update("status", todo.Status)
	}
	return tx.Commit().Error
}

func DeleteTodos(ids []uint) error {
	var todos []models.Todo
	db.StorageDB.Where("id in (?)", ids).Find(&todos)
	return db.StorageDB.Delete(&todos).Error
}

func CreateProject(project models.Project) error {
	err := db.StorageDB.Create(&project).Error
	return err
}

func ListProjects(userId string) ([]ListProjectResponseItem, error) {
	var projects []models.Project
	db.StorageDB.Where("user_id = ?", userId).Find(&projects)

	var response []ListProjectResponseItem

	for _, project := range projects {
		lri := ListProjectResponseItem{
			ID:          project.ID,
			Title:       project.Title,
			Description: project.Description,
			ProjectId:   project.ID,
			UserId:      project.UserID.String(),
		}
		response = append(response, lri)
	}
	return response, nil
}

func DeleteProject(id uint) error {
	var project models.Project
	db.StorageDB.Where("id = ?", id).Find(&project)
	return db.StorageDB.Delete(&project).Error
}
