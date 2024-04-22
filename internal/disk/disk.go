package disk

import (
	"archive/tar"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/lpoirothattermann/storage/internal/bundler"
)

func WriteBundleToPath(bundleReader *bundler.BundleReader, directoryOutput string) {
	for true {
		header, err := bundleReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error while scanning archive: %v\n", err)
		}

		filepath := filepath.Join(directoryOutput, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(filepath, 0755); err != nil {
				log.Fatalf("Error while creating folder: %v", err)
			}
		case tar.TypeReg:
			if err := WriteBufferToFilePath(directoryOutput, header.Name, bundleReader.TarReader); err != nil {
				log.Fatalf("Error while writing file on disk: %v\n", err)
			}
		default:
			log.Fatalf(
				"Unknown type: %v in %v",
				header.Typeflag,
				header.Name,
			)
		}
	}
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
		log.Fatalf("Error while reading path: %v\n", err)
	}

	return true
}

func CreateSymlink(source string, destination string, force bool) {
	if FileOrDirectoryExists(destination) {
		if force == false {
			log.Fatalf("Symlink already exists")
		}
		if err := os.Remove(destination); err != nil {
			log.Fatalf("Error while deleting old symlink: %v\n", err)
		}
	}

	if err := os.Symlink(source, destination); err != nil {
		log.Fatalf("Error while creating symlink: %v\n", err)
	}
}
