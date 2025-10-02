package server

import (
	"archivus-v2/internal/service"
	"archivus-v2/pkg/logging"
	reqhelpers "archivus-v2/pkg/reqHelpers"
	"archivus-v2/pkg/response"
	"bytes"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func GetFilesByFolder(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")
	folder := r.URL.Query().Get("folder")
	files, folderSize, err := service.GetFiles(userId, folder)
	if err != nil {
		logging.Errorlogger.Error().Msg(err.Error())
		response.InternalServerErrorResponse(w, err.Error())
		return
	}
	response.JSONResponse(w, map[string]interface{}{
		"files": files,
		"size":  folderSize,
	})
}

func GetSignedUrlHandler(w http.ResponseWriter, r *http.Request) {
	filepath := mux.Vars(r)["filepath"]
	userId := r.Header.Get("userId")
	signedUrl, err := service.GetSignedUrl(filepath, userId)
	// ownership check please
	if err != nil {
		logging.Errorlogger.Error().Msg(err.Error())
		response.InternalServerErrorResponse(w, err.Error())
		return
	}
	response.JSONResponse(w, map[string]interface{}{
		"signed_url": signedUrl})
}

func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	compressed := false
	filepath := mux.Vars(r)["filepath"]
	signature := r.URL.Query().Get("signature")
	expiresAtStr := r.URL.Query().Get("expires_at")
	compressedStr := r.URL.Query().Get("compressed")
	if compressedStr == "true" {
		compressed = true
	}
	expiresAt, err := strconv.Atoi(expiresAtStr)
	if err != nil {
		logging.Errorlogger.Error().Msg(err.Error())
		response.BadRequestResponse(w, err.Error())
		return
	}
	if time.Now().After(time.Unix(int64(expiresAt), 0)) {
		response.UnauthorizedResponse(w, "Signature expired")
		return
	}
	f, err := service.DownloadFile(filepath, signature, expiresAtStr, compressed)
	if err != nil {
		logging.Errorlogger.Error().Msg(err.Error())
		response.InternalServerErrorResponse(w, err.Error())
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath)
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeContent(w, r, filepath, time.Now(), bytes.NewReader(f))
}

type MoveFileRequest struct {
	FilePath string `json:"filePath"`
	Dst      string `json:"dst"`
}

func MoveFileHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")
	var req MoveFileRequest
	err := reqhelpers.DecodeRequest(r, &req)
	if err != nil {
		logging.Errorlogger.Error().Msg(err.Error())
		response.BadRequestResponse(w, "Invalid request body")
		return
	}
	err = service.MoveFile(userId, req.FilePath, req.Dst, false)
	if err != nil {
		logging.Errorlogger.Error().Msg(err.Error())
		response.InternalServerErrorResponse(w, err.Error())
		return
	}
	response.JSONResponse(w, map[string]interface{}{"message": "File moved successfully"})
}

type DeleteFileRequest struct {
	FilePath string `json:"filePath"`
}

func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")
	var req DeleteFileRequest
	err := reqhelpers.DecodeRequest(r, &req)
	if err != nil {
		logging.Errorlogger.Error().Msg(err.Error())
		response.BadRequestResponse(w, "Invalid request body")
		return
	}

	err = service.DeleteFile(userId, req.FilePath)
	if err != nil {
		logging.Errorlogger.Error().Msg(err.Error())
		response.InternalServerErrorResponse(w, err.Error())
		return
	}
	response.JSONResponse(w, map[string]interface{}{"message": "File deleted successfully"})
}
