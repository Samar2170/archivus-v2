package config

import (
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
	UploadsDir         string               `yaml:"uploads_dir"`
	BackendConfig      HostingConfiguration `yaml:"backend_config"`
	FrontEndConfig     HostingConfiguration `yaml:"frontend_config"`
}

var Config *Configuration

var OS string
var Arch string

func runTimeInfo() {
	OS = runtime.GOOS
	Arch = runtime.GOARCH
}

func LoadConfig(filePath string) error {
	// homeDir, homeDirErr := os.UserHomeDir()
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, &Config); err != nil {
		return err
	}
	// if Config == nil {
	// 	Config = &Configuration{}
	// }
	// if Config.LogsDir == "" {
	// 	Config.LogsDir = "logs"
	// }
	// if Config.StorageDbFile == "" {
	// 	Config.StorageDbFile = "storage.db"
	// }
	// if Config.UploadsDir != "" && !filepath.IsAbs(Config.UploadsDir) {
	// 	Config.UploadsDir = filepath.Join(BaseDir, Config.UploadsDir)
	// }
	// if Config.UploadsDir == "" {
	// 	if homeDirErr != nil {
	// 		return errors.New("failed to get user home directory")
	// 	} else {
	// 		Config.UploadsDir = filepath.Join(homeDir)
	// 	}
	// }
	// if Config.MasterPinUploads && Config.MasterPin == "" {
	// 	return errors.New("masterPin is required when MasterPinUploads is true")
	// }
	return nil
}

var BaseDir string

const (
	HeaderName         = "AccessKey"
	CompressionQuality = 85
)

func init() {
	var err error
	currentFile, err := os.Executable()
	if err != nil {
		panic(err)
	}

	BaseDir = filepath.Dir(filepath.Dir(currentFile))
	// err = LoadConfig(filepath.Join(BaseDir, "config.yaml"))
	err = LoadConfig("../config.yaml")
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}
	CheckConfig()
	runTimeInfo()
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
		Config.FrontEndConfig.BaseUrl == "" ||
		Config.FrontEndConfig.Port == "" {
		panic("Backend or Frontend configuration is incomplete")
	}

}
