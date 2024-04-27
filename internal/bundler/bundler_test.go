package bundler_test

import (
	"bytes"
	"testing"

	"filippo.io/age"
	"github.com/lpoirothattermann/storage/internal/bundler"
	"github.com/lpoirothattermann/storage/internal/constants"
	"github.com/stretchr/testify/assert"
)

func TestBundleWriterClose(t *testing.T) {
	ageIdentity, err := age.GenerateX25519Identity()
	assert.Nil(t, err)

	buffer := bytes.Buffer{}
	bundleWriter, err := bundler.NewWriter(&buffer, ageIdentity.Recipient())
	assert.Nil(t, err)
	err = bundleWriter.Close()
	assert.Nil(t, err)
}

func TestGetFinalFilename(t *testing.T) {
	testList := []struct {
		filenameWithoutExtension string
		expectedResult           string
	}{
		{"secret", "secret" + constants.ARCHIVE_EXTENSION},
	}

	for _, test := range testList {
		assert.Equal(t, test.expectedResult, bundler.GetFinalFilename(test.filenameWithoutExtension))
	}
}

func TestNewWriter(t *testing.T) {
	buffer := bytes.Buffer{}

	ageIdentity, err := age.GenerateX25519Identity()
	assert.Nil(t, err)

	_, err = bundler.NewWriter(&buffer, ageIdentity.Recipient())
	assert.Nil(t, err)
}

func TestNewReader(t *testing.T) {
	ageIdentity, err := age.GenerateX25519Identity()
	assert.Nil(t, err)

	buffer := bytes.Buffer{}
	bundleWriter, err := bundler.NewWriter(&buffer, ageIdentity.Recipient())
	assert.Nil(t, err)
	err = bundleWriter.Close()
	assert.Nil(t, err)

	_, err = bundler.NewReader(&buffer, ageIdentity)
	assert.Nil(t, err)
}
