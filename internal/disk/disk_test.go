package disk_test

import (
	"archive/tar"
	"bytes"
	"io/fs"
	"os"
	"path"
	"testing"

	"filippo.io/age"
	"github.com/lpoirothattermann/storage/internal/bundler"
	"github.com/lpoirothattermann/storage/internal/disk"
	"github.com/stretchr/testify/assert"
)

func TestFileOrDirectoryExistsTrue(t *testing.T) {
	tmpDir := t.TempDir()

	testList := []struct {
		_type   string
		dirPath string
		name    string
	}{
		{"dir", "path/to", "test"},
		{"file", "path/to", "test"},
		{"symlink", "path/to", "test"},
	}

	for _, test := range testList {
		var pathToTest string
		baseTmpDir, err := os.MkdirTemp(tmpDir, "*")
		assert.Nil(t, err)

		switch test._type {
		case "dir":
			pathToTest = path.Join(baseTmpDir, test.dirPath, test.name)
			err = os.MkdirAll(pathToTest, fs.ModePerm)
			assert.Nil(t, err)
		case "file":
			baseTmpDir = path.Join(baseTmpDir, test.dirPath)
			err = os.MkdirAll(baseTmpDir, fs.ModePerm)
			assert.Nil(t, err)
			pathToTest = path.Join(baseTmpDir, test.name)
			_, err = os.Create(pathToTest)
			assert.Nil(t, err)
		case "symlink":
			pathDirSourceForSymlink, err := os.MkdirTemp(tmpDir, "*")
			assert.Nil(t, err)

			baseTmpDir = path.Join(baseTmpDir, test.dirPath)
			err = os.MkdirAll(baseTmpDir, fs.ModePerm)
			assert.Nil(t, err)

			pathToTest = path.Join(baseTmpDir, test.name)
			err = disk.CreateSymlink(pathDirSourceForSymlink, pathToTest, true)
			assert.Nil(t, err)
		}

		assert.True(t, disk.FileOrDirectoryExists(pathToTest))
	}
}

func TestFileOrDirectoryExistsFalse(t *testing.T) {
	assert.False(t, disk.FileOrDirectoryExists("/path/that/do/not/exists"))
}

func TestCreateSymlinkValid(t *testing.T) {
	tmpDir := t.TempDir()

	testList := []struct {
		fileAlreadyExists bool
	}{
		{true},
		{false},
	}

	for _, test := range testList {
		pathDirSourceForSymlink, err := os.MkdirTemp(tmpDir, "*")
		assert.Nil(t, err)

		symlinkPathToTest := path.Join(tmpDir, "symlinkTarget")

		if test.fileAlreadyExists {
			_, err = os.Create(symlinkPathToTest)
			assert.Nil(t, err)
			assert.Nil(t, disk.CreateSymlink(pathDirSourceForSymlink, symlinkPathToTest, true))
		} else {
			assert.Nil(t, disk.CreateSymlink(pathDirSourceForSymlink, symlinkPathToTest, false))
		}

		assert.True(t, disk.FileOrDirectoryExists(symlinkPathToTest))

		err = os.RemoveAll(symlinkPathToTest)
		assert.Nil(t, err)
	}
}

func TestCreateSymlinkInvalid(t *testing.T) {
	tmpDir := t.TempDir()

	pathDirSourceForSymlink, err := os.MkdirTemp(tmpDir, "*")
	assert.Nil(t, err)

	symlinkPathToTest := path.Join(tmpDir, "symlinkTarget")

	_, err = os.Create(symlinkPathToTest)
	assert.Nil(t, err)
	assert.NotNil(t, disk.CreateSymlink(pathDirSourceForSymlink, symlinkPathToTest, false))
	assert.NotNil(t, disk.CreateSymlink(pathDirSourceForSymlink, path.Join(symlinkPathToTest, "/path/that/not/exits"), true))
	assert.NotNil(t, disk.CreateSymlink(path.Join(pathDirSourceForSymlink, "/path/that/not/exits"), symlinkPathToTest, true))
}

func TestWriteBufferToFilePath(t *testing.T) {
	tmpDir := t.TempDir()

	const stringToTest = "Secure and Protected"
	bytesBuffer := bytes.Buffer{}
	bytesBuffer.Write([]byte(stringToTest))

	const filename = "test.txt"
	filePathToTest := path.Join(tmpDir, filename)
	disk.WriteBufferToFilePath(tmpDir, filename, &bytesBuffer)

	bytesContent, err := os.ReadFile(filePathToTest)
	assert.Nil(t, err)
	content := string(bytesContent)

	assert.Contains(t, content, stringToTest)
}

func TestBundleAndUnbundleToPath(t *testing.T) {
	// Create bundle
	tmpDir := t.TempDir()
	folderName := "secret"
	outputFilename := bundler.GetFinalFilename(folderName)

	ageIdentity, err := age.GenerateX25519Identity()
	assert.Nil(t, err)

	tarballBuffer := bytes.Buffer{}

	bundleWriter, err := bundler.NewWriter(&tarballBuffer, ageIdentity.Recipient())
	assert.Nil(t, err)

	tarHeader := &tar.Header{
		Name:     folderName,
		Mode:     0755,
		Typeflag: tar.TypeDir,
	}

	err = bundleWriter.TarWriter.WriteHeader(tarHeader)
	assert.Nil(t, err)
	err = bundleWriter.Close()
	assert.Nil(t, err)

	err = disk.WriteBufferToFilePath(tmpDir, outputFilename, &tarballBuffer)
	assert.Nil(t, err)

	bundlePath := path.Join(tmpDir, outputFilename)
	bytesContent, err := os.ReadFile(bundlePath)
	assert.Nil(t, err)
	content := string(bytesContent)

	assert.Contains(t, content, "age-encryption.org/v1")
	assert.Contains(t, content, "X25519")

	// Unbundle
	file, err := os.Open(bundlePath)
	assert.Nil(t, err)
	bundleReader, err := bundler.NewReader(file, ageIdentity)
	assert.Nil(t, err)
	err = disk.UnbundleToPath(bundleReader, tmpDir)
	assert.Nil(t, err)
	assert.True(t, disk.FileOrDirectoryExists(path.Join(tmpDir, folderName)))
}
