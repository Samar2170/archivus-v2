package server

import (
	"archivus-v2/config"
	"archivus-v2/internal"
	"archivus-v2/internal/middleware"
	"archivus-v2/pkg/logging"
	"archivus-v2/pkg/response"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response.SuccessResponse(w, "OK")
}

func GetServer(testEnv bool) *http.Server {
	internal.SetupRun(testEnv)

	logger := logging.AuditLogger
	mux := mux.NewRouter()

	mux.HandleFunc("/health", HealthCheck)
	mux.HandleFunc("/login", Login)

	mux.HandleFunc("/files/upload/", UploadFilesHandler).Methods("POST")

	mux.HandleFunc("/files/get/", GetFilesByFolder)
	mux.HandleFunc("/files/get-signed-url/{filepath:.*}", GetSignedUrlHandler)
	mux.HandleFunc("/files/download/{filepath:.*}", DownloadFileHandler)

	mux.HandleFunc("/folder/add/", CreateFolderHandler).Methods("POST")

	mux.HandleFunc("/files/move/", MoveFileHandler).Methods("POST")
	mux.HandleFunc("/files/delete/", DeleteFileHandler).Methods("POST")

	mux.HandleFunc("/bigupload/initiate/", InitiateBigUpload).Methods("POST")
	mux.HandleFunc("/bigupload/chunk/", UploadChunk).Methods("POST")
	mux.HandleFunc("/bigupload/finalize/", FinalizeBigUpload).Methods("POST")

	logMiddleware := logging.NewLogMiddleware(&logger)
	mux.Use(logMiddleware.Func())
	wrappedMux := middleware.AuthMiddleware(mux)
	wrappedMux = CorsConfig.Handler(wrappedMux)

	server := &http.Server{
		Handler: wrappedMux,
		Addr:    config.GetBackendAddr(),
	}
	return server
}

func RunServer() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	server := GetServer(false)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Errorlogger.Printf("Failed to start server: %v", err)
		}
	}()
	fmt.Printf("Server is running at %s\n", config.GetBackendAddr())
	<-stop

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server gracefully stopped")
}
