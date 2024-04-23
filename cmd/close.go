package cmd

import (
	"archive/tar"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/lpoirothattermann/storage/internal/bundler"
	"github.com/lpoirothattermann/storage/internal/config"
	"github.com/lpoirothattermann/storage/internal/disk"
	"github.com/lpoirothattermann/storage/internal/log"
	"github.com/spf13/cobra"
)

func init() {
	command := &cobra.Command{
		Use:   "close [state]",
		Short: "Close state",
		Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run:   closeCmdFunc,
	}

	rootCmd.AddCommand(command)
}

func closeCmdFunc(cmd *cobra.Command, args []string) {
	stateName := args[0]

	state, exists := config.GetConfig().States[stateName]
	if exists == false {
		log.Critical.Fatalf("State %q doesn't exists.\n'", stateName)
	}

	if disk.FileOrDirectoryExists(state.GetTemporaryDirectoryPath()) == false {
		log.Critical.Fatalf("State %q can't be close becose it's not open.\n", stateName)
	}

	// Backup archive
	if err := os.Rename(state.GetArchivePath(), state.GetArchiveBackupPath()); err != nil {
		log.Critical.Fatalf("Error while making backup of your archive: %v", err)
	}

	// Create empty archive to receive all data later
	archiveFile, err := os.Create(state.GetArchivePath())
	if err != nil {
		log.Critical.Fatalf("Error while creating empty archive: %v\n", err)
	}

	bundleWriter, err := bundler.NewWriter(archiveFile, state.AgeIdentity.Recipient())
	if err != nil {
		log.Critical.Fatalf("Error while creating bundle writer: %v\n", err)
	}

	err = filepath.Walk(state.GetTemporaryDirectoryPath(), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(state.GetTemporaryDirectoryPath(), path)
		if err != nil {
			return err
		}
		header.Name = relPath

		if err := bundleWriter.TarWriter.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			// Copy the file data into the tarball
			if _, err := io.Copy(bundleWriter.TarWriter, file); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		log.Critical.Fatalf("Error while coping data in the archive: %v\n", err)
	}

	if err := bundleWriter.Close(); err != nil {
		log.Critical.Fatalf("Error while closing bundle writer: %v\n", err)
	}

	if err := archiveFile.Close(); err != nil {
		log.Critical.Fatalf("Error while closing archive: %v\n", err)
	}

	if err := os.RemoveAll(state.GetSymlinkTargetPath()); err != nil {
		log.Critical.Fatalf("Error while removing symlink %q -> %q: %v\n", state.GetTemporaryDirectoryPath(), state.GetSymlinkTargetPath(), err)
	}

	if err := os.RemoveAll(state.GetTemporaryDirectoryPath()); err != nil {
		log.Critical.Fatalf("Error while removing %q: %v\n", state.GetTemporaryDirectoryPath(), err)
	}

	if err := os.RemoveAll(state.GetArchiveBackupPath()); err != nil {
		log.Critical.Fatalf("Error while removing archive backup %q: %v\n", state.GetArchiveBackupPath(), err)
	}
}
