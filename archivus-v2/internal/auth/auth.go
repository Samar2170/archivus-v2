package auth

import (
	"archivus-v2/config"
	"archivus-v2/internal/db"
	dirmanager "archivus-v2/internal/dirManager"
	"archivus-v2/internal/models"
	"archivus-v2/internal/utils"
	"errors"

	"github.com/google/uuid"
)

func CreateUser(username, password, pin, email string, createDir, isMaster bool) (string, string, error) {
	// Validate input
	if username == "" || pin == "" || password == "" {
		return "", "", utils.HandleError("CreateUser", "Username, PIN, and email cannot be empty", errors.New("invalid input"))
	}

	// Generate API key
	apiKey, err := utils.GenerateAPIKey(config.ApiKeyLength)
	if err != nil {
		return "", "", utils.HandleError("CreateUser", "Failed to generate API key", err)
	}

	// Hash the PIN
	hashedPin := utils.HashString(pin)
	hashedPassword := utils.HashString(password)
	tx := db.StorageDB.Begin()
	user := models.User{
		ID:       uuid.New(),
		Username: username,
		PIN:      hashedPin,
		Email:    email,
		APIKey:   apiKey,
		Password: hashedPassword,
		IsMaster: isMaster,
	}
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return "", "", utils.HandleError("CreateUser", "Failed to create user in database", err)
	}
	if config.Config.UserDirStructure && createDir {
		err = dirmanager.CreateDirForUser(user)
		if err != nil {
			tx.Rollback()
			return "", "", utils.HandleError("CreateUser", "Failed to create user folder", err)
		}
	}
	if err := tx.Commit().Error; err != nil {
		return "", "", utils.HandleError("CreateUser", "Failed to commit transaction", err)
	}

	return apiKey, user.ID.String(), nil
}
