package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

type HostingConfiguration struct {
	BaseUrl  string `yaml:"base_url"`
	Port     string `yaml:"port"`
	Scheme   string `yaml:"scheme"`
	BindAddr string `yaml:"bind_addr"`
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
	AllowedOrigins     []string             `yaml:"allowed_origins"`
	ServerSalt         string               `yaml:"server_salt"`
	BackendProxyUrl    string               `yaml:"backend_proxy_url"`
}

var Config *Configuration

func (c *Configuration) Print() {
	fmt.Println("Configuration:")
	fmt.Printf("  Mode: %s\n", c.Mode)
	fmt.Printf("  MultiUser: %v\n", c.MultiUser)
	fmt.Printf("  DefaultWriteAccess: %v\n", c.DefaultWriteAccess)
	fmt.Printf("  UserDirStructure: %v\n", c.UserDirStructure)
	fmt.Printf("  UserDirLock: %v\n", c.UserDirLock)
	fmt.Printf("  MasterPinUploads: %v\n", c.MasterPinUploads)
	fmt.Printf("  MasterPin: %s\n", c.MasterPin)
	fmt.Printf("  LogsDir: %s\n", c.LogsDir)
	fmt.Printf("  StorageDbFile: %s\n", c.StorageDbFile)
	fmt.Printf("  SecretKey: %s\n", c.SecretKey)
	fmt.Printf("  BotToken: %s\n", c.BotToken)
	fmt.Printf("  BaseDir: %s\n", c.BaseDir)
	fmt.Printf("  ServerSalt: %s\n", c.ServerSalt)
	fmt.Printf("  AllowedOrigins: %v\n", c.AllowedOrigins)

	fmt.Println("  BackendConfig:")
	fmt.Printf("    BaseUrl: %s\n", c.BackendConfig.BaseUrl)
	fmt.Printf("    Port: %s\n", c.BackendConfig.Port)
	fmt.Printf("    Scheme: %s\n", c.BackendConfig.Scheme)
	fmt.Printf("    BindAddr: %s\n", c.BackendConfig.BindAddr)

	fmt.Println("  FrontEndConfig:")
	fmt.Printf("    BaseUrl: %s\n", c.FrontEndConfig.BaseUrl)
	fmt.Printf("    Port: %s\n", c.FrontEndConfig.Port)
	fmt.Printf("    Scheme: %s\n", c.FrontEndConfig.Scheme)
	fmt.Printf("    BindAddr: %s\n", c.FrontEndConfig.BindAddr)
}

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
