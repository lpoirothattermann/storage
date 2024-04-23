package age

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"filippo.io/age"
)

func GetIdentityFromFile(filePath string) (*age.X25519Identity, error) {
	privateKeyFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer privateKeyFile.Close()

	fileScanner := bufio.NewScanner(privateKeyFile)

	var privateKeyString string
	for fileScanner.Scan() {
		if strings.HasPrefix(fileScanner.Text(), "AGE-SECRET-KEY-1") {
			privateKeyString = fileScanner.Text()
			break
		}
	}

	if privateKeyString == "" {
		return nil, errors.New("No private key found in file")
	}

	identity, err := age.ParseX25519Identity(privateKeyString)
	if err != nil {
		return nil, err
	}

	return identity, nil
}
