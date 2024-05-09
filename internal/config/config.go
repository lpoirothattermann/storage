package config

import (
	"log"
	"os"
	"path"

	"github.com/lpoirothattermann/storage/internal/constants"
	logInternal "github.com/lpoirothattermann/storage/internal/log"
	pathInternal "github.com/lpoirothattermann/storage/internal/path"
	"github.com/spf13/viper"
)

type rawConfig struct {
	States      map[string]rawState `mapstructure:"states"`
	LogFilePath string              `mapstructure:"log_file_path"`
}

type Config struct {
	LogFilePath string
}

var (
	cacheConfig    *Config
	cacheRawConfig *rawConfig
	cacheStates    map[string]*State
)

func GetConfig() *Config {
	return cacheConfig
}

func init() {
	cacheStates = make(map[string]*State)
	cacheRawConfig = &rawConfig{}
	cacheConfig = getConfig()
}

func getConfig() *Config {
	configDirectoryPath := getConfigDirectoryPath()

	// Ensure that config directory is created
	if err := os.MkdirAll(configDirectoryPath, os.ModePerm); err != nil {
		format := "Error while creating config directory: %v\n"
		if logInternal.Critical != nil {
			logInternal.Critical.Fatalf(format, err)
		} else {
			log.Fatalf(format, err)
		}
	}

	viper.AddConfigPath(configDirectoryPath)
	viper.SetConfigName(constants.CONFIG_BASENAME)
	viper.SetConfigType(constants.CONFIG_FILETYPE)

	// Load config
	// TODO not a better way to do?
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			createDefaultConfigFile()
			getConfig()
		} else {
			logInternal.Critical.Fatalf("Error while loading config: %v\n", err)
		}
	}
	if err := viper.UnmarshalExact(cacheRawConfig); err != nil {
		logInternal.Critical.Fatalf("Error while parsing configuration file: %v\n", err)
	}

	// Create formatted config struct with empty states
	return &Config{
		LogFilePath: pathInternal.GetNormalizedPath(cacheRawConfig.LogFilePath),
	}
}

func (config *Config) GetState(stateName string) *State {
	state, exists := cacheStates[stateName]
	if exists {
		return state
	}

	rawState, exists := cacheRawConfig.States[stateName]
	if !exists {
		logInternal.Critical.Fatalf("State %q doesn't exists.", stateName)
	}

	cacheStates[stateName] = &State{
		Name:          stateName,
		ageIdentity:   nil,
		ageRecipient:  nil,
		EncryptedPath: pathInternal.GetNormalizedPath(rawState.EncryptedPath),
		DecryptedPath: pathInternal.GetNormalizedPath(rawState.DecryptedPath),
	}

	return cacheStates[stateName]
}

func getConfigDirectoryPath() string {
	// TODO improve, process in constants file directly ?
	if os.Getenv("GOENV") == "dev" {
		return "test_data/"
	}

	return constants.CONFIG_DIRECTORY_PATH
}

func GetConfigFilePath() string {
	// TODO process in constants file directly ?
	return path.Join(getConfigDirectoryPath(), constants.CONFIG_BASENAME+"."+constants.CONFIG_FILETYPE)
}

func createDefaultConfigFile() {
	setDefaultConfigs()
	if err := viper.WriteConfigAs(GetConfigFilePath()); err != nil {
		logInternal.Critical.Fatalf("Error while creating config file %q: %v\n", GetConfigFilePath(), err)
	}
}

// Set every default configurations, it will be used to create default configuration file
func setDefaultConfigs() {
	viper.Set("log_file_path", path.Join(getConfigDirectoryPath(), constants.CONFIG_DEFAULT_LOG_FILE_FILENAME))
	viper.Set("states", make(map[string]struct{}, 0))
}
