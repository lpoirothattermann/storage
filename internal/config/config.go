package config

import (
	"log"
	"os"

	ageInternal "github.com/lpoirothattermann/storage/internal/age"
	"github.com/lpoirothattermann/storage/internal/constants"
	"github.com/lpoirothattermann/storage/internal/format"
	"github.com/lpoirothattermann/storage/internal/path"
	"github.com/spf13/viper"
)

type rawConfig struct {
	States map[string]rawState `mapstructure:"states"`
}

type Config struct {
	States map[string]State
}

func GetConfig() *Config {
	configDirectoryPath := getConfigDirectoryPath()

	// Ensure that config directory is created
	if err := os.MkdirAll(configDirectoryPath, os.ModePerm); err != nil {
		log.Fatalf("Error while creating config directory: %v\n", err)
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
			log.Fatal(format.ErrorTypeAndMessage("main", err))
		}
	}

	rawConfig := &rawConfig{}
	if err := viper.UnmarshalExact(rawConfig); err != nil {
		format.ErrorTypeAndMessage("config", err)
	}

	// Create formated config struct
	config := Config{
		States: make(map[string]State),
	}

	for index, value := range rawConfig.States {
		identity := ageInternal.GetIdentityFromFile(path.GetNormalizedPath(value.PrivateKeyPath))

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
		log.Fatal(format.ErrorTypeAndMessage("main", err))
	}
}

// Set every default configurations, it will be used to create default configuration file
func setDefaultConfigs() {
	viper.Set("states", make(map[string]struct{}, 0))
}
