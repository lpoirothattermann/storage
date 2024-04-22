package cmd

import (
	"archive/tar"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/lpoirothattermann/storage/internal/bundler"
	"github.com/lpoirothattermann/storage/internal/config"
	"github.com/lpoirothattermann/storage/internal/disk"
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
		log.Fatalf("State %q doesn't exists.\n'", stateName)
	}

	// Backup old encrypt path
	os.Rename(state.GetArchivePath(), state.GetArchiveBackupPath())

	// Bundle tmp state to encrypt path
	archiveFile, err := os.Create(state.GetArchivePath())
	if err != nil {
		log.Fatalf("Error while creating empty archive: %v\n", err)
	}

	if disk.FileOrDirectoryExists(state.GetTemporaryDirectoryPath()) == false {
		log.Fatalf("State is not open\n")
	}

	bundleWriter, err := bundler.NewWriter(archiveFile, state.AgeIdentity.Recipient())
	if err != nil {
		log.Fatalf("Error while creating bundle writer: %v\n", err)
	}

	filepath.Walk(state.GetTemporaryDirectoryPath(), func(path string, info fs.FileInfo, err error) error {
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

	bundleWriter.Close()
	archiveFile.Close()

	os.RemoveAll(state.GetSymlinkTargetPath())
	os.RemoveAll(state.GetTemporaryDirectoryPath())
	os.RemoveAll(state.GetArchiveBackupPath())
}
