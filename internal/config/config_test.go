package config_test

import (
	"path"
	"testing"

	"github.com/lpoirothattermann/storage/internal/config"
	"github.com/lpoirothattermann/storage/internal/constants"
	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	// How to test without risking personal config file intergrity?
}

func TestReloadConfig(t *testing.T) {
	// How to test without risking personal config file intergrity?
}

func TestGetConfigFilePath(t *testing.T) {
	listTest := []struct {
		ValueGOENV     string
		ExpectedOutput string
	}{
		{"dev", path.Join(constants.CONFIG_RELATIVE_DEV_DIRECTORY_PATH, constants.CONFIG_BASENAME+"."+constants.CONFIG_FILETYPE)},
		{"hello", path.Join(constants.CONFIG_DIRECTORY_PATH, constants.CONFIG_BASENAME+"."+constants.CONFIG_FILETYPE)},
		{"", path.Join(constants.CONFIG_DIRECTORY_PATH, constants.CONFIG_BASENAME+"."+constants.CONFIG_FILETYPE)},
	}

	for _, test := range listTest {
		t.Setenv("GOENV", test.ValueGOENV)
		assert.Equal(t, test.ExpectedOutput, config.GetConfigFilePath())
	}
}
