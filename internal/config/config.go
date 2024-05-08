package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"filippo.io/age"
	ageInternal "github.com/lpoirothattermann/storage/internal/age"
	"github.com/lpoirothattermann/storage/internal/constants"
	logInternal "github.com/lpoirothattermann/storage/internal/log"
	pathInternal "github.com/lpoirothattermann/storage/internal/path"
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
			return GetConfig()
		} else {
			logInternal.Critical.Fatalf("Error while loading config: %v\n", err)
		}
	}

	rawConfig := &rawConfig{}
	if err := viper.UnmarshalExact(rawConfig); err != nil {
		logInternal.Critical.Fatalf("Error while parsing configuration file: %v\n", err)
	}

	// Create formated config struct
	config := Config{
		LogFilePath: pathInternal.GetNormalizedPath(rawConfig.LogFilePath),
		States:      make(map[string]State),
	}

	for index, value := range rawConfig.States {
		privateKeyPath := pathInternal.GetNormalizedPath(value.PrivateKeyPath)

		privateKeyReader, err := os.Open(privateKeyPath)
		if err != nil {
			logInternal.Critical.Fatalf("Unable to open private key file: %v\n", err)
		}
		privateKeyReaderLimited := io.LimitReader(privateKeyReader, constants.PRIVATE_KEY_SIZE_LIMIT)

		if ageInternal.IsEncryptedWithPassphrase(privateKeyReaderLimited) {
			privateKeyReader.Seek(0, io.SeekStart)

			passphrase, err := ageInternal.PassphrasePromptForDecryption()
			if err != nil {
				logInternal.Critical.Fatalf("Error while getting passphrase from user input: %v\n", err)
			}

			scryptIdentity, err := age.NewScryptIdentity(passphrase)
			if err != nil {
				logInternal.Critical.Fatalf("Error while creating passphrase identity: %v\n", err)
			}
			ad, err := age.Decrypt(privateKeyReaderLimited, scryptIdentity)
			if err != nil {
				if _, isWrongPassphrase := err.(*age.NoIdentityMatchError); isWrongPassphrase {
					fmt.Printf("Wrong passphrase.\n")
					os.Exit(1)
				} else {
					logInternal.Critical.Fatalf("Error while decrypting key with passphrase: %T\n", err)
				}
			}

			privateKeyFileDecrypted, err := io.ReadAll(ad)
			if err != nil {
				logInternal.Critical.Fatalf("Error while parsing key file: %v\n", err)
			}
			privateKeyReaderLimited = strings.NewReader(string(privateKeyFileDecrypted))
		}

		identities, err := age.ParseIdentities(privateKeyReaderLimited)
		if err != nil {
			logInternal.Critical.Fatalf("Error while getting identity from key file: %v\n", err)
		}

		config.States[index] = State{
			Name:          index,
			AgeIdentity:   identities[0].(*age.X25519Identity),
			EncryptedPath: pathInternal.GetNormalizedPath(value.EncryptedPath),
			DecryptedPath: pathInternal.GetNormalizedPath(value.DecryptedPath),
		}
	}

	return &config
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
