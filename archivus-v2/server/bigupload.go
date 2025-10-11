package server

import (
	"archivus-v2/internal/service/biguploads"
	"net/http"
)

func InitiateBigUpload(w http.ResponseWriter, r *http.Request) {
	biguploads.InitBigUpload(w, r)
}

func UploadChunk(w http.ResponseWriter, r *http.Request) {
	biguploads.ChunkHandler(w, r)
}

func FinalizeBigUpload(w http.ResponseWriter, r *http.Request) {
	biguploads.FinaliseHandler(w, r)
}
