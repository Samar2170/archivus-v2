package config

import "path/filepath"

const (
	TestDbFile        = "test.db"
	ConfigFile        = "config.yaml"
	UserId            = "userId"
	Username          = "username"
	ApiKeyLength      = 32
	PasswordMinLength = 8
	PinLelength       = 6
	CredsFile         = "creds.json"
	DefaultUserDir    = "users"
	UserInfoFileName  = ".userinfo.json"

	sessionsDir = ".sessions"
	TmpSuffix   = ".tmp"
	MaxChunks   = 1024
)

func GetSessionsDir() string {
	return filepath.Join(Config.BaseDir, sessionsDir)
}

func GetTmpDir() string {
	return filepath.Join(Config.BaseDir, TmpSuffix)
}
