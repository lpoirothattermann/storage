package age_test

import (
	"os"
	"path"
	"testing"

	"filippo.io/age"
	ageInternal "github.com/lpoirothattermann/storage/internal/age"
	"github.com/stretchr/testify/assert"
)

func TestGetIdentityFromFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create key.txt
	ageIdentity, err := age.GenerateX25519Identity()
	assert.Nil(t, err)

	ageIdentityKeyPath := path.Join(tmpDir, "age.txt")
	file, err := os.Create(ageIdentityKeyPath)

	_, err = file.WriteString(ageIdentity.String())
	assert.Nil(t, err)
	err = file.Close()
	assert.Nil(t, err)

	// ---

	_, err = ageInternal.GetIdentityFromFile(ageIdentityKeyPath)
	assert.Nil(t, err)
}
