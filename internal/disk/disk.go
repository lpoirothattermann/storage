package disk

import (
	"archive/tar"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/lpoirothattermann/storage/internal/bundler"
	"github.com/lpoirothattermann/storage/internal/log"
)

func UnbundleToPath(bundleReader *bundler.BundleReader, directoryOutput string) error {
	for true {
		header, err := bundleReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		filepath := filepath.Join(directoryOutput, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(filepath, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := WriteBufferToFilePath(directoryOutput, header.Name, bundleReader.TarReader); err != nil {
				return err
			}
		default:
			log.Critical.Fatalf(
				"Unknown type: %v in %v",
				header.Typeflag,
				header.Name,
			)
		}
	}

	return nil
}

func WriteBufferToFilePath(path string, fileName string, buffer io.Reader) error {
	file, err := os.Create(filepath.Join(path, fileName))
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, buffer); err != nil {
		return err
	}

	return nil
}

func FileOrDirectoryExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else if err != nil {
		return false
	}

	return true
}

func CreateSymlink(source string, destination string, force bool) error {
	if !FileOrDirectoryExists(source) {
		return errors.New("Source directory doesn't exists")
	}

	if FileOrDirectoryExists(destination) {
		if force == false {
			return errors.New("Symlink already exists")
		}
		if err := os.RemoveAll(destination); err != nil {
			return err
		}
	}

	if err := os.Symlink(source, destination); err != nil {
		return err
	}

	return nil
}
