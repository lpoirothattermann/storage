package config

import (
	"os"

	ageInternal "github.com/lpoirothattermann/storage/internal/age"
	"github.com/lpoirothattermann/storage/internal/constants"
	"github.com/lpoirothattermann/storage/internal/log"
	"github.com/lpoirothattermann/storage/internal/path"
	"github.com/spf13/viper"
)

type rawConfig struct {
	LogFilePath string              `mapstructure:"log_file_path"`
	States      map[string]rawState `mapstructure:"states"`
}

type Config struct {
	LogFilePath string
	States      map[string]State
}

// Singleton
var singleConfig *Config

func GetConfig() *Config {
	if singleConfig == nil {
		singleConfig = getConfig()
	}

	return singleConfig
}

func ReloadConfig() *Config {
	singleConfig = getConfig()

	return singleConfig
}

func getConfig() *Config {
	configDirectoryPath := getConfigDirectoryPath()

	// Ensure that config directory is created
	if err := os.MkdirAll(configDirectoryPath, os.ModePerm); err != nil {
		log.Critical.Fatalf("Error while creating config directory: %v\n", err)
	}

	viper.AddConfigPath(configDirectoryPath)
	viper.SetConfigName(constants.CONFIG_BASENAME)
	viper.SetConfigType(constants.CONFIG_FILETYPE)

	// Load config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			createDefaultConfigFile()
			return GetConfig()
		} else {
			log.Critical.Fatalf("Error while loading config: %v\n", err)
		}
	}

	rawConfig := &rawConfig{}
	if err := viper.UnmarshalExact(rawConfig); err != nil {
		log.Critical.Fatalf("Error while parsing configuration file: %v\n", err)
	}

	// Create formated config struct
	config := Config{
		LogFilePath: path.GetNormalizedPath(rawConfig.LogFilePath),
		States:      make(map[string]State),
	}

	for index, value := range rawConfig.States {
		privateKeyPath := path.GetNormalizedPath(value.PrivateKeyPath)

		identity, err := ageInternal.GetIdentityFromFile(privateKeyPath)
		if err != nil {
			log.Critical.Fatalf("Error while getting identity from file %q: %v\n", privateKeyPath, err)
		}

		config.States[index] = State{
			Name:          index,
			AgeIdentity:   identity,
			EncryptedPath: path.GetNormalizedPath(value.EncryptedPath),
			DecryptedPath: path.GetNormalizedPath(value.DecryptedPath),
		}
	}

	return &config
}

func getConfigDirectoryPath() string {
	// TODO improve
	if os.Getenv("GOENV") == "dev" {
		return "test_data/"
	}

	return constants.CONFIG_DIRECTORY_PATH
}

func GetConfigFilePath() string {
	return getConfigDirectoryPath() + constants.CONFIG_BASENAME + "." + constants.CONFIG_FILETYPE
}

func createDefaultConfigFile() {
	setDefaultConfigs()
	if err := viper.WriteConfigAs(GetConfigFilePath()); err != nil {
		log.Critical.Fatalf("Error while creating config file %q: %v\n", GetConfigFilePath(), err)
	}
}

// Set every default configurations, it will be used to create default configuration file
func setDefaultConfigs() {
	viper.Set("states", make(map[string]struct{}, 0))
}
