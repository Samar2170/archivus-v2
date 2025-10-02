package server

import (
	"archivus-v2/internal/auth"
	reqhelpers "archivus-v2/pkg/reqHelpers"
	"archivus-v2/pkg/response"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	// Parse the login request
	var loginReq auth.LoginUserRequest
	err := reqhelpers.DecodeRequest(r, &loginReq)
	if err != nil {
		response.BadRequestResponse(w, "Invalid request body")
		return
	}

	// Validate the login credentials
	token, userId, err := auth.LoginUser(loginReq)
	if err != nil {
		response.UnauthorizedResponse(w, err.Error())
		return
	}

	// Generate a session token for the user
	response.JSONResponse(w, map[string]interface{}{
		"token": token, "user_id": userId,
	})

}
