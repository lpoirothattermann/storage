package cmd

import (
	"archive/tar"
	"bytes"
	"io"
	"os"

	"filippo.io/age"
	ageInternal "github.com/lpoirothattermann/storage/internal/age"
	"github.com/lpoirothattermann/storage/internal/bundler"
	"github.com/lpoirothattermann/storage/internal/constants"
	"github.com/lpoirothattermann/storage/internal/disk"
	logInternal "github.com/lpoirothattermann/storage/internal/log"
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
	stateName := args[2]

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
		logInternal.Critical.Fatalf("Error while getting identity from file %q: %v\n", privateKeyPath, err)
	}
	ageIdentity := identities[0].(*age.X25519Identity)

	tarballBuffer := bytes.Buffer{}

	bundleWriter, err := bundler.NewWriter(&tarballBuffer, ageIdentity.Recipient())
	if err != nil {
		logInternal.Critical.Fatalf("Error while openning bundle writer: %v\n", err)
	}

	tarHeader := &tar.Header{
		Name:     stateName,
		Mode:     0o755,
		Typeflag: tar.TypeDir,
	}

	if err := bundleWriter.TarWriter.WriteHeader(tarHeader); err != nil {
		logInternal.Critical.Fatalf("Error while writing header to tarball: %v\n", err)
	}

	if err := bundleWriter.Close(); err != nil {
		logInternal.Critical.Fatalf("Error closing bundle writer: %v\n", err)
	}

	if err := disk.WriteBufferToFilePath(stateDirectoryOutput, bundler.GetFinalFilename(stateName), &tarballBuffer); err != nil {
		logInternal.Critical.Fatalf("Error while writing tarball buffer to file: %v\n", err)
	}

	logInternal.Info.Printf("Archive succesfully created, archive can be used for a new state with name %q.\n", stateName)
}
