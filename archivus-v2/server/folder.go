package server

import (
	"archivus-v2/internal/dirmanager"
	"archivus-v2/pkg/logging"
	reqhelpers "archivus-v2/pkg/reqHelpers"
	"archivus-v2/pkg/response"
	"net/http"
)

func CreateFolderHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
	var err error
	userWriteAccess, err := checkUserWriteAccess(username)
	if err != nil {
		logging.Errorlogger.Error().Msg(err.Error())
		response.InternalServerErrorResponse(w, err.Error())
		return
	}
	if !userWriteAccess {
		logging.Errorlogger.Error().Msg("User does not have write access")
		response.ForbiddenResponse(w, "User does not have write access")
		return
	}

	type CreateFolderRequest struct {
		Folder string
	}
	var req CreateFolderRequest
	err = reqhelpers.DecodeRequest(r, &req)
	if err != nil {
		logging.Errorlogger.Error().Msg(err.Error())
		response.BadRequestResponse(w, "Invalid request body")
		return
	}
	if req.Folder == "" {
		response.BadRequestResponse(w, "Folder path is required")
		return
	}
	err = dirmanager.CreateFolder(req.Folder, username)
	if err != nil {
		logging.Errorlogger.Error().Msg(err.Error())
		response.InternalServerErrorResponse(w, err.Error())
		return
	}
	response.SuccessResponse(w, "Folder created successfully")
}
