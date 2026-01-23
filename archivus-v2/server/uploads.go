package server

import (
	"archivus-v2/internal/service"
	"archivus-v2/pkg/logging"
	"archivus-v2/pkg/response"
	"net/http"
)

func UploadFilesHandler(w http.ResponseWriter, r *http.Request) {
	// Set max request body size to 500 MB for video support
	r.Body = http.MaxBytesReader(w, r.Body, 500<<20)

	// Parse multipart form with 256 MB max memory (rest goes to temp files)
	err := r.ParseMultipartForm(256 << 20)
	if err != nil {
		logging.Errorlogger.Error().Msg(err.Error())
		response.BadRequestResponse(w, err.Error())
		return
	}

	username := r.Header.Get("username")
	folderPath := r.FormValue("folder")

	// Get multiple files with field name "file"
	files := r.MultipartForm.File["file"]
	if len(files) == 0 {
		response.BadRequestResponse(w, "No files provided")
		return
	}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			logging.Errorlogger.Error().Msg(err.Error())
			response.BadRequestResponse(w, "Error opening file: "+fileHeader.Filename)
			return
		}
		defer file.Close()

		err = service.SaveFile(file, fileHeader, username, folderPath, "")
		if err != nil {
			logging.Errorlogger.Error().Msg(err.Error())
			response.InternalServerErrorResponse(w, "Error saving file: "+fileHeader.Filename)
			return
		}
	}
	response.SuccessResponse(w, "File uploaded successfully")
}
