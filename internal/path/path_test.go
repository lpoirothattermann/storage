package path_test

import (
	"os"
	"testing"

	"github.com/lpoirothattermann/storage/internal/path"
	"github.com/stretchr/testify/assert"
)

func TestGetNormalizedPath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	assert.Nil(t, err)

	testList := []struct {
		path           string
		expectedResult string
	}{
		{"~/.config/storage", homeDir + "/.config/storage"},
		{"~/.config/storage/", homeDir + "/.config/storage"},
		{"~/storage", homeDir + "/storage"},
		{"~/storage/", homeDir + "/storage"},
		{"/tmp/storage", "/tmp/storage"},
		{"/tmp/storage/", "/tmp/storage"},
	}

	for _, test := range testList {
		assert.Equal(
			t,
			test.expectedResult,
			path.GetNormalizedPath(test.path),
		)
	}
}
