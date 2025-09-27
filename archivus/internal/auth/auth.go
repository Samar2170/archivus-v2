package auth

import (
	"archivus/config"
	"archivus/internal/db"
	"archivus/internal/helpers"
	"archivus/internal/models"
	"archivus/internal/utils"

	"github.com/google/uuid"
)

func CreateUser(username, password, pin, email string, isMaster bool) (string, string, error) {
	// Validate input
	if username == "" || pin == "" || email == "" || password == "" {
		return "", "", utils.HandleError("CreateUser", "Username, PIN, and email cannot be empty", nil)
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
	if config.Config.UserDirStructure {
		err = helpers.CreateFolderForUser(user.Username)
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

type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	PIN      string `json:"pin"`
}

func LoginUser(req LoginUserRequest) (string, string, error) {
	var user models.User
	var err error
	var token string
	var userId string
	if req.Username == "" || (req.Password == "" && req.PIN == "") {
		return token, userId, utils.HandleError("LoginUser", "Username, password, and PIN cannot be empty", nil)
	}
	if req.PIN == "" {
		err = db.StorageDB.Where("username = ?", req.Username).
			Where("password = ?", utils.HashString(req.Password)).
			First(&user).Error
	} else {
		err = db.StorageDB.Where("username = ?", req.Username).
			Where("pin = ?", utils.HashString(req.PIN)).
			First(&user).Error
	}
	userId = user.ID.String()
	if err != nil {
		return token, userId, utils.HandleError("LoginUser", "Invalid credentials", err)
	}
	token, err = createToken(user.ID.String(), user.Username)
	if err != nil {
		return token, userId, utils.HandleError("LoginUser", "Failed to create token", err)
	}
	return token, userId, nil
}

func GetUserByApiKey(apiKey string) (models.User, error) {
	var user models.User
	err := db.StorageDB.Where("api_key = ?", apiKey).First(&user).Error
	return user, err
}

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
	PIN      string `json:"pin"`
	ApiKey   string `json:"api_key"`
}

// func CreateUserByMasterPerm(userFolder string) (Credential, error) {
// 	var cd Credential
// 	cd.Username = userFolder
// 	apiKey, err := utils.GenerateAPIKey(config.ApiKeyLength)
// 	if err != nil {
// 		logging.Errorlogger.Error().Msgf("Error generating API key: %v", err)
// 		fmt.Printf("Error generating %s API key: %s", userFolder, err)
// 		return cd, err
// 	}
// 	cd.ApiKey = apiKey
// 	password, err := utils.GenerateAPIKey(config.PasswordMinLength)
// 	if err != nil {
// 		logging.Errorlogger.Error().Msgf("Error generating password: %v", err)
// 		fmt.Printf("Error generating %s password: %s", userFolder, err)
// 		return cd, err
// 	}
// 	cd.Password = password
// 	pin, err := utils.GenerateRandomNumber(config.PinLelength)
// 	if err != nil {
// 		logging.Errorlogger.Error().Msgf("Error generating PIN: %v", err)
// 		fmt.Printf("Error generating %s PIN: %s", userFolder, err)
// 		return cd, err
// 	}
// 	cd.PIN = pin
// 	tx := db.StorageDB.Begin()
// 	var existingUser models.User
// 	err = db.StorageDB.Where("username = ?", userFolder).First(&existingUser).Error
// 	if err != gorm.ErrRecordNotFound {
// 		tx.Rollback()
// 		logging.Errorlogger.Error().Msgf("Error creating user for %s: %v", userFolder, err)
// 		fmt.Printf("Error searching for existing user for %s: %s\n", userFolder, err)
// 		return cd, err
// 	}
// 	tx.Create(&models.User{
// 		ID:       uuid.New(),
// 		Username: userFolder,
// 		Email:    userFolder + "@placeholder.com",
// 		PIN:      utils.HashString(pin),
// 		Password: utils.HashString(password),
// 		APIKey:   apiKey,
// 	})
// 	err = tx.Commit().Error
// 	if err != nil {
// 		logging.Errorlogger.Error().Msgf("Error creating user for %s: %v", userFolder, err)
// 		fmt.Printf("Error creating user for %s: %s\n", userFolder, err)
// 		return cd, err
// 	} else {
// 		fmt.Printf("User created for %s with API key: %s\n", userFolder, apiKey)
// 		return cd, err
// 	}
// }
