package log_test

import (
	"os"
	"path"
	"testing"

	"github.com/lpoirothattermann/storage/internal/constants"
	"github.com/lpoirothattermann/storage/internal/log"
	"github.com/stretchr/testify/assert"
)

func TestInitializationWithCreatedLogFile(t *testing.T) {
	tmpDir := t.TempDir()

	createLogFile := func(dirPath string, filename string) string {
		err := os.MkdirAll(tmpDir, os.ModePerm)
		assert.Nil(t, err)
		dirPath, err = os.MkdirTemp(tmpDir, "*")
		assert.Nil(t, err)

		logFilePath := path.Join(dirPath, filename)
		_, err = os.Create(logFilePath)
		assert.Nil(t, err)

		return logFilePath
	}

	testList := []struct {
		logFileDirectoryPath string
		logFileFilename      string
	}{
		{path.Join(tmpDir, "/my_super_dir"), "storage"},
		{path.Join(tmpDir, "/my_beautiful_dir"), "storage.log"},
	}

	for _, test := range testList {
		logFilePath := createLogFile(test.logFileDirectoryPath, test.logFileFilename)

		assert.Nil(t, log.Initialization(logFilePath))

		messageLogPartCritical := "Bonjour"
		messageLogPartError := "Hello"
		messageLogPartWarning := "こんにちは"
		messageLogPartInfo := "Здравствуйте"

		log.Critical.Println(messageLogPartCritical)
		log.Error.Println(messageLogPartError)
		log.Warning.Println(messageLogPartWarning)
		log.Info.Println(messageLogPartInfo)

		assert.Nil(t, log.Close())

		bytesContent, err := os.ReadFile(logFilePath)
		assert.Nil(t, err)
		content := string(bytesContent)

		assert.Contains(t, content, messageLogPartCritical)
		assert.Contains(t, content, messageLogPartError)
		assert.Contains(t, content, messageLogPartWarning)
		assert.Contains(t, content, messageLogPartInfo)
		assert.Contains(t, content, constants.LOG_PREFIX_CRITICAL)
		assert.Contains(t, content, constants.LOG_PREFIX_ERROR)
		assert.Contains(t, content, constants.LOG_PREFIX_WARNING)
		assert.Contains(t, content, constants.LOG_PREFIX_INFO)
	}
}

func TestInitializationWithoutCreatedLogFile(t *testing.T) {
	logFilePathThatNotExists := "/storage/is/a/great/tool/storage.log"

	assert.NotNil(t, log.Initialization(logFilePathThatNotExists))
	assert.NotNil(t, log.Close())
}
