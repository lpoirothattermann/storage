package cmd

import (
	"archive/tar"
	"bytes"
	"fmt"
	"log"

	"filippo.io/age"
	ageInternal "github.com/lpoirothattermann/storage/internal/age"
	"github.com/lpoirothattermann/storage/internal/bundler"
	"github.com/lpoirothattermann/storage/internal/disk"
	"github.com/lpoirothattermann/storage/internal/format"
	"github.com/lpoirothattermann/storage/internal/path"
	"github.com/spf13/cobra"
)

func init() {
	command := &cobra.Command{
		Use:   "init [private key path] [state directory output] [state name]",
		Short: "Init a new state by creating new .tar.gz.age archive that contains only base directory, ready to be use as a new state archive",
		Args:  cobra.MatchAll(cobra.ExactArgs(3), cobra.OnlyValidArgs),
		Run:   initCmdFunc,
	}
	rootCmd.AddCommand(command)
}

func initCmdFunc(cmd *cobra.Command, args []string) {
	privateKeyPath := path.GetNormalizedPath(args[0])
	stateDirectoryOutput := path.GetNormalizedPath(args[1])
	nameWithoutExtensions := args[2]

	ageIdentity := ageInternal.GetIdentityFromFile(privateKeyPath)

	if err := createArchiveForNewState(ageIdentity.Recipient(), stateDirectoryOutput, nameWithoutExtensions); err != nil {
		log.Fatal(format.ErrorTypeAndMessage("cmd", err))
	}

	fmt.Println("Archive succesfully created, it can be use for a new state.")
}

func createArchiveForNewState(recipient *age.X25519Recipient, directoryPath string, stateName string) error {
	archiveBaseDirectoryName := stateName
	tarballBuffer := bytes.Buffer{}

	bundleWriter, err := bundler.NewWriter(&tarballBuffer, recipient)
	if err != nil {
		fmt.Println("Error while openning bundler writer:", err)

		return err
	}

	tarHeader := &tar.Header{
		Name:     fmt.Sprintf("%v/", archiveBaseDirectoryName),
		Mode:     0755,
		Typeflag: tar.TypeDir,
	}

	if err := bundleWriter.TarWriter.WriteHeader(tarHeader); err != nil {
		fmt.Println("Error writing directory header to tarball:", err)

		return err
	}

	if err := bundleWriter.Close(); err != nil {
		fmt.Println("Error closing bundle writer:", err)

		return err
	}

	if err := disk.WriteBufferToFilePath(directoryPath, bundler.GetFinalFilename(stateName), &tarballBuffer); err != nil {
		fmt.Println("Error while writing tarball to file:", err)

		return err
	}

	return nil
}
