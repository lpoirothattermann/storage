package config

import (
	"io"
	"os"
	"path/filepath"

	"filippo.io/age"
	ageInternal "github.com/lpoirothattermann/storage/internal/age"
	"github.com/lpoirothattermann/storage/internal/constants"
	"github.com/lpoirothattermann/storage/internal/disk"
	"github.com/lpoirothattermann/storage/internal/log"
	"github.com/lpoirothattermann/storage/internal/path"
)

type rawState struct {
	PrivateKeyPath string `mapstructure:"private_key_path"`
	EncryptedPath  string `mapstructure:"encrypted_path"`
	DecryptedPath  string `mapstructure:"decrypted_path"`
	PublicKey      string `mapstructure:"public_key"`
}

type State struct {
	Name          string
	ageIdentity   *age.X25519Identity  // Cache
	ageRecipient  *age.X25519Recipient // Cache
	EncryptedPath string
	DecryptedPath string
}

func (state *State) GetIdentity() *age.X25519Identity {
	if state.ageIdentity == nil {
		state.ageIdentity = state.getIdentity()
	}

	return state.ageIdentity
}

func (state *State) getIdentity() *age.X25519Identity {
	rawState := cacheRawConfig.States[state.Name]
	privateKeyPath := path.GetNormalizedPath(rawState.PrivateKeyPath)

	privateKeyReader, err := os.Open(privateKeyPath)
	if err != nil {
		log.Critical.Fatalf("Unable to open private key file: %v\n", err)
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
		log.Critical.Fatalf("Error while getting identity from key file %q: %v\n", privateKeyPath, err)
	}

	return identities[0].(*age.X25519Identity)
}

func (state *State) GetRecipient() *age.X25519Recipient {
	if state.ageRecipient == nil {
		if rawState := cacheRawConfig.States[state.Name]; rawState.PublicKey == "" {
			state.ageRecipient = state.GetIdentity().Recipient()
		} else {
			recipient, err := age.ParseX25519Recipient(rawState.PublicKey)
			if err != nil {
				log.Critical.Fatalf("Error while parsing provided public key: %v\n", err)
			}
			state.ageRecipient = recipient
		}
	}

	return state.ageRecipient
}

func (state *State) IsOpen() bool {
	return disk.FileOrDirectoryExists(state.GetTemporaryDirectoryPath())
}

func (state *State) IsClose() bool {
	return !state.IsOpen()
}

func (state *State) GetTemporaryDirectoryPath() string {
	return filepath.Join(os.TempDir(), constants.APPLICATION_TEMPORARY_DIRECTORY_NAME, state.Name)
}

func (state *State) GetSymlinkTargetPath() string {
	return filepath.Join(state.DecryptedPath, state.Name)
}

func (state *State) GetArchivePath() string {
	return filepath.Join(state.EncryptedPath, state.Name+constants.ARCHIVE_EXTENSION)
}

func (state *State) GetArchiveBackupPath() string {
	return state.GetArchivePath() + constants.ARCHIVE_BACKUP_EXTENSION
}
