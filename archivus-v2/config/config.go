package config

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

type HostingConfiguration struct {
	BaseUrl string `yaml:"base_url"`
	Port    string `yaml:"port"`
	Scheme  string `yaml:"scheme"`
}

type Configuration struct {
	Mode               string               `yaml:"mode"`
	MultiUser          bool                 `yaml:"multi_user"`
	DefaultWriteAccess bool                 `yaml:"default_write_access"`
	UserDirStructure   bool                 `yaml:"user_dir_structure"`
	UserDirLock        bool                 `yaml:"user_dir_lock"`
	MasterPinUploads   bool                 `yaml:"master_pin_uploads"`
	MasterPin          string               `yaml:"master_pin"`
	LogsDir            string               `yaml:"logs_dir"`
	StorageDbFile      string               `yaml:"storage_db_file"`
	SecretKey          string               `yaml:"secret_key"`
	BotToken           string               `yaml:"bot_token"` // Added BotToken to the configuration
	BaseDir            string               `yaml:"base_dir"`
	BackendConfig      HostingConfiguration `yaml:"backend_config"`
	FrontEndConfig     HostingConfiguration `yaml:"frontend_config"`
	ServerSalt         string               `yaml:"server_salt"`
}

var Config *Configuration

var OS string
var Arch string

var ProjectBaseDir string

const (
	HeaderName         = "AccessKey"
	CompressionQuality = 85
)

func SetupConfig(mode string) {
	var err error
	var defaultConfigPaths []string

	if mode == "prod" {
		currentFile, err := os.Executable()
		if err != nil {
			panic(err)
		}
		ProjectBaseDir = filepath.Dir(currentFile)

		defaultConfigPaths = []string{
			ProdConfigFile,
			"../" + ProdConfigFile,
			filepath.Join(ProjectBaseDir, ProdConfigFile),
		}
	} else {
		ProjectBaseDir = "."
		defaultConfigPaths = []string{
			DevConfigFile,
			"../" + DevConfigFile,
			filepath.Join(ProjectBaseDir, DevConfigFile),
		}
	}

	for _, path := range defaultConfigPaths {
		err = LoadConfig(path)
		if err == nil {
			break
		}
	}
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}
	CheckConfig()
	runTimeInfo()
}

func runTimeInfo() {
	OS = runtime.GOOS
	Arch = runtime.GOARCH
}

func LoadConfig(filePath string) error {
	homeDir, homeDirErr := os.UserHomeDir()
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, &Config); err != nil {
		return err
	}
	if Config == nil {
		log.Fatal("Configuration is nil after unmarshalling")
	}
	if Config.LogsDir == "" {
		Config.LogsDir = "logs"
	}
	if Config.StorageDbFile == "" {
		Config.StorageDbFile = "storage.db"
	}
	if Config.BaseDir == "" {
		if homeDirErr != nil {
			return errors.New("failed to get user home directory")
		} else {
			Config.BaseDir = filepath.Join(homeDir)
		}
	}
	if Config.MasterPinUploads && Config.MasterPin == "" {
		return errors.New("masterPin is required when MasterPinUploads is true")
	}
	if Config.ServerSalt == "" {
		salt, err := createOrLoadServerSalt()
		if err != nil {
			return err
		}
		Config.ServerSalt = salt
	}
	return nil
}

func CheckConfig() {
	if Config == nil {
		panic("Configuration is nil after loading")
	}
	if Config.Mode == "" {
		panic("Mode is not set in the configuration")
	}
	if Config.SecretKey == "" {
		panic("SecretKey is not set in the configuration")
	}
	if Config.BackendConfig.BaseUrl == "" ||
		Config.BackendConfig.Port == "" ||
		Config.FrontEndConfig.Port == "" {
		panic("Backend or Frontend configuration is incomplete")
	}

}
