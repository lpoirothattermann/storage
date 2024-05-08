package config

import (
	"os"
	"path/filepath"

	"filippo.io/age"
	"github.com/lpoirothattermann/storage/internal/constants"
	"github.com/lpoirothattermann/storage/internal/disk"
)

type rawState struct {
	PrivateKeyPath string `mapstructure:"private_key_path"`
	EncryptedPath  string `mapstructure:"encrypted_path"`
	DecryptedPath  string `mapstructure:"decrypted_path"`
}

type State struct {
	Name          string
	AgeIdentity   *age.X25519Identity
	EncryptedPath string
	DecryptedPath string
}

func (state *State) IsOpen() bool {
	if !disk.FileOrDirectoryExists(state.GetTemporaryDirectoryPath()) {
		return false
	}

	return true
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
