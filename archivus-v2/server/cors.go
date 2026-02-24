package server

import (
	"archivus-v2/config"
	"archivus-v2/pkg/logging"
	"fmt"

	"github.com/rs/cors"
)

var CorsConfig *cors.Cors

func SetupCors() {
	allowedOrigins := config.Config.AllowedOrigins

	// Add default local frontend URL if not present
	defaultFrontend := fmt.Sprintf("http://%s:%s", config.Config.FrontEndConfig.BaseUrl, config.Config.FrontEndConfig.Port)
	for i := 1; i < 255; i++ {
		allowedOrigins = append(allowedOrigins, fmt.Sprintf("http://192.168.1.%d:%s", i, config.Config.FrontEndConfig.Port))
	}
	for i := 1; i < 10; i++ {
		allowedOrigins = append(allowedOrigins, fmt.Sprintf("http://localhost:%d", 3000+i))
	}
	allowedOrigins = append(allowedOrigins, "http://localhost:3000")

	found := false
	for _, origin := range allowedOrigins {
		if origin == defaultFrontend {
			found = true
			break
		}
	}
	if !found {
		allowedOrigins = append(allowedOrigins, defaultFrontend)
	}

	logger := cors.Logger(&logging.AuditLogger)

	CorsConfig = cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE", "PATCH", "HEAD"},
		AllowedHeaders:   []string{"Origin", "Accept", "Content-Type", "Authorization", "X-Requested-With", "AccessKey"},
		AllowCredentials: true,
		Logger:           logger,
	})

}
