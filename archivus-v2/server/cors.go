package server

import (
	"archivus-v2/config"
	"archivus-v2/pkg/logging"
	"fmt"

	"github.com/rs/cors"
)

var CorsConfig *cors.Cors

func SetupCors() {
	allowedOrigins := []string{}
	for i := 1; i < 255; i++ {
		allowedOrigins = append(allowedOrigins, fmt.Sprintf("http://192.168.1.%d:%s", i, config.Config.FrontEndConfig.Port))
	}
	for i := 1; i < 10; i++ {
		allowedOrigins = append(allowedOrigins, fmt.Sprintf("http://localhost:%d", 3000+i))
	}
	allowedOrigins = append(allowedOrigins, "http://localhost:3000")
	allowedOrigins = append(allowedOrigins, fmt.Sprintf("http://%s:%s", config.Config.FrontEndConfig.BaseUrl, config.Config.FrontEndConfig.Port))
	allowedOrigins = append(allowedOrigins, fmt.Sprintf("http://%s:%s/", config.Config.FrontEndConfig.BaseUrl, config.Config.FrontEndConfig.Port))
	allowedOrigins = append(allowedOrigins, fmt.Sprintf("%s://%s", config.Config.BackendNetworkConfg.Scheme, config.Config.BackendNetworkConfg.BaseUrl))
	logger := cors.Logger(&logging.AuditLogger)

	CorsConfig = cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE", "PATCH", "HEAD"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Logger:           logger,
	})

}
