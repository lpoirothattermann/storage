package age

import (
	"bufio"
	"log"
	"os"
	"strings"

	"filippo.io/age"
)

func GetIdentityFromFile(filePath string) *age.X25519Identity {
	privateKeyFile, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error while opening private key file: %v\n", err)
	}

	fileScanner := bufio.NewScanner(privateKeyFile)

	var privateKeyString string
	for fileScanner.Scan() {
		if strings.HasPrefix(fileScanner.Text(), "AGE-SECRET-KEY-1") {
			privateKeyString = fileScanner.Text()
			break
		}
	}
	privateKeyFile.Close()

	if privateKeyString == "" {
		log.Fatal("No private key in the given file")
	}

	identity, err := age.ParseX25519Identity(privateKeyString)
	if err != nil {
		log.Fatalf("Error while parsing private key: %v\n", err)
	}

	return identity
}
