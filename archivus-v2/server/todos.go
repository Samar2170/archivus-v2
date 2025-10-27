package server

import (
	"archivus-v2/internal/models"
	"archivus-v2/internal/tempora"
	"archivus-v2/pkg/response"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func Todos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var todos []tempora.TodoRequest
		userId := r.Header.Get(UserIdKey)
		err := json.NewDecoder(r.Body).Decode(&todos)
		if err != nil {
			response.BadRequestResponse(w, err.Error())
			return
		}
		err = tempora.CreateTodos(todos, userId)
		if err != nil {
			response.InternalServerErrorResponse(w, err.Error())
			return
		}
		response.SuccessResponse(w, "Todos created successfully")
		return
	case http.MethodGet:
		userId := r.Header.Get("userId")
		projectId := r.URL.Query().Get("projectId")
		todos, err := tempora.ListTodos(userId, projectId)
		if err != nil {
			response.InternalServerErrorResponse(w, err.Error())
			return
		}
		response.JSONResponse(w, todos)
		return
	}
}

func UpdateTodos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var updateTodos []tempora.UpdateTodosRequest
		err := json.NewDecoder(r.Body).Decode(&updateTodos)
		if err != nil {
			response.BadRequestResponse(w, err.Error())
			return
		}
		err = tempora.UpdateTodos(updateTodos)
		if err != nil {
			response.InternalServerErrorResponse(w, err.Error())
			return
		}
		response.SuccessResponse(w, "Todos marked as done successfully")
	case http.MethodDelete:
		var ids []uint
		err := json.NewDecoder(r.Body).Decode(&ids)
		if err != nil {
			response.BadRequestResponse(w, err.Error())
			return
		}
		err = tempora.DeleteTodos(ids)
		if err != nil {
			response.InternalServerErrorResponse(w, err.Error())
			return
		}
		response.SuccessResponse(w, "Todos deleted successfully")
	}
}
func parseUUID(s string) (uuid.UUID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, err
	}
	return u, nil
}

func Projects(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		userId := r.Header.Get(UserIdKey)
		projects, err := tempora.ListProjects(userId)
		if err != nil {
			response.InternalServerErrorResponse(w, err.Error())
			return
		}
		response.JSONResponse(w, projects)
		return
	case http.MethodPost:
		userId := r.Header.Get(UserIdKey)
		var project models.Project
		err := json.NewDecoder(r.Body).Decode(&project)
		if err != nil {
			response.BadRequestResponse(w, err.Error())
			return
		}
		project.UserID, err = parseUUID(userId)
		if err != nil {
			response.InternalServerErrorResponse(w, err.Error())
			return
		}
		err = tempora.CreateProject(project)
		if err != nil {
			response.InternalServerErrorResponse(w, err.Error())
			return
		}
		response.SuccessResponse(w, "Project created successfully")
		return
	case http.MethodDelete:
		var id uint
		err := json.NewDecoder(r.Body).Decode(&id)
		if err != nil {
			response.BadRequestResponse(w, err.Error())
			return
		}
		err = tempora.DeleteProject(id)
		if err != nil {
			response.InternalServerErrorResponse(w, err.Error())
			return
		}
		response.SuccessResponse(w, "Project deleted successfully")
		return
	}
}
