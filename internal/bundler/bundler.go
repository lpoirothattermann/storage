package bundler

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"

	"filippo.io/age"
	"github.com/lpoirothattermann/storage/internal/config"
	"github.com/lpoirothattermann/storage/internal/constants"
)

type BundleWriter struct {
	TarWriter  *tar.Writer
	GzipWriter *gzip.Writer
	AgeWriter  io.WriteCloser
}
type BundleReader struct {
	TarReader  *tar.Reader
	GzipReader io.Reader
	AgeReader  io.Reader
}

func (bundleWriter *BundleWriter) Close() error {
	if err := bundleWriter.TarWriter.Close(); err != nil {
		fmt.Println("Error closing tar writer:", err)

		return err
	}

	if err := bundleWriter.GzipWriter.Close(); err != nil {
		fmt.Println("Error closing gzip writer:", err)

		return err
	}

	if err := bundleWriter.AgeWriter.Close(); err != nil {
		fmt.Println("Error closing age writer:", err)

		return err
	}

	return nil
}

func (bundleReader *BundleReader) Next() (*tar.Header, error) {
	return bundleReader.TarReader.Next()
}

func (bundleReader *BundleReader) Read(bytes []byte) (int, error) {
	return bundleReader.TarReader.Read(bytes)
}

func GetFinalFilenameFromState(state config.State) string {
	return GetFinalFilename(state.Name)
}

func GetFinalFilename(filenameWitoutExtension string) string {
	return fmt.Sprintf("%v%v", filenameWitoutExtension, constants.ARCHIVE_EXTENSION)
}

func NewWriter(buffer io.Writer, ageRecipient *age.X25519Recipient) (*BundleWriter, error) {
	ageWriter, err := age.Encrypt(buffer, ageRecipient)
	if err != nil {
		fmt.Println("Error opening age writer", err)

		return nil, err
	}

	gzipWriter := gzip.NewWriter(ageWriter)
	tarWriter := tar.NewWriter(gzipWriter)

	return &BundleWriter{
		TarWriter:  tarWriter,
		GzipWriter: gzipWriter,
		AgeWriter:  ageWriter,
	}, nil
}

func NewReader(buffer io.Reader, ageIdentity *age.X25519Identity) (*BundleReader, error) {
	ageReader, err := age.Decrypt(buffer, ageIdentity)
	if err != nil {
		return nil, err
	}

	gzipReader, err := gzip.NewReader(ageReader)
	if err != nil {
		return nil, err
	}

	tarReader := tar.NewReader(gzipReader)

	return &BundleReader{
		TarReader:  tarReader,
		GzipReader: gzipReader,
		AgeReader:  ageReader,
	}, nil
}
