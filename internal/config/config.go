package config

import (
	"io"
	"log"
	"os"
	"path"

	"filippo.io/age"
	ageInternal "github.com/lpoirothattermann/storage/internal/age"
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

func (config *Config) GetState(stateName string) *State {
	if _, exists := statesCache[stateName]; !exists {
		statesCache[stateName] = getState(stateName)
	}

	return statesCache[stateName]
}

// Singleton
var (
	vConfig     *Config
	vRawConfig  *rawConfig
	statesCache map[string]*State
)

func GetConfig() *Config {
	return vConfig
}

func init() {
	vRawConfig = &rawConfig{}
	vConfig = getConfig()
	statesCache = make(map[string]*State)
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
	if err := viper.UnmarshalExact(vRawConfig); err != nil {
		logInternal.Critical.Fatalf("Error while parsing configuration file: %v\n", err)
	}

	// Create formatted config struct with empty states
	return &Config{
		LogFilePath: pathInternal.GetNormalizedPath(vRawConfig.LogFilePath),
	}
}

func getState(stateName string) *State {
	for sName, sValue := range vRawConfig.States {
		if sName != stateName {
			continue
		}

		privateKeyPath := pathInternal.GetNormalizedPath(sValue.PrivateKeyPath)

		privateKeyReader, err := os.Open(privateKeyPath)
		if err != nil {
			logInternal.Critical.Fatalf("Unable to open private key file: %v\n", err)
		}
		privateKeyReaderLimited := io.LimitReader(privateKeyReader, constants.PRIVATE_KEY_SIZE_LIMIT)

		if ageInternal.IsEncryptedWithPassphrase(privateKeyReaderLimited) {
			privateKeyReader.Seek(0, io.SeekStart)
			privateKeyReaderLimited = ageInternal.AskAndDecryptWithPassphrase(privateKeyReaderLimited)
		} else {
			privateKeyReader.Seek(0, io.SeekStart)
		}

		identities, err := age.ParseIdentities(privateKeyReaderLimited)
		if err != nil {
			logInternal.Critical.Fatalf("Error while getting identity from key file %q: %v\n", privateKeyPath, err)
		}

		return &State{
			Name:          sName,
			AgeIdentity:   identities[0].(*age.X25519Identity),
			EncryptedPath: pathInternal.GetNormalizedPath(sValue.EncryptedPath),
			DecryptedPath: pathInternal.GetNormalizedPath(sValue.DecryptedPath),
		}
	}

	logInternal.Critical.Fatalf("State %q doesn't exists.", stateName)

	return nil // Not reachable
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
