package middleware

import (
	"net/http"
	"strings"

	"archivus-v2/config"
	"archivus-v2/internal/auth"
	"archivus-v2/internal/models"
	"archivus-v2/pkg/response"
)

var ExemptPaths = map[string]struct{}{"/files/download/": {}, "/login": {}}

func CheckExemptPath(path string) bool {
	for exemptPath := range ExemptPaths {
		if len(path) > len(exemptPath) {
			if path[:len(exemptPath)] == exemptPath {
				return true
			}
		} else {
			if path == exemptPath {
				return true
			}
		}
	}
	return false
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Preflight requests must bypass auth — they carry no credentials
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}
		exempt := CheckExemptPath(r.URL.Path)
		if exempt {
			next.ServeHTTP(w, r)
			return
		}
		tokenString := r.Header.Get("Authorization")
		token := ""
		if tokenString != "" {
			parts := strings.Split(tokenString, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			} else {
				response.BadRequestResponse(w, "Invalid Authorization header format")
				return
			}
		}
		apiKeyHeader := r.Header.Get("X-API-Key")
		if token == "" && apiKeyHeader == "" {
			response.UnauthorizedResponse(w, "Missing JWT Token or API Key")
			return
		}
		var username, userId string
		var err error
		if token != "" {
			userId, username, err = auth.DecodeToken(token)
			if err != nil {
				response.UnauthorizedResponse(w, "Invalid JWT Token")
				return
			}
		} else {
			user, err := models.GetUserByApiKey(apiKeyHeader)
			if err != nil {
				response.UnauthorizedResponse(w, "Invalid API Key")
				return
			}
			userId = user.ID.String()
			username = user.Username
		}
		r.Header.Set(config.Username, username)
		r.Header.Set(config.UserId, userId)
		next.ServeHTTP(w, r)
	})
}

// handle jwt token logic
