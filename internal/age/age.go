package age

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"filippo.io/age"
	logInternal "github.com/lpoirothattermann/storage/internal/log"
	"github.com/lpoirothattermann/storage/internal/tui"
)

func IsEncryptedWithPassphrase(reader io.Reader) bool {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "-> scrypt") {
			return true
		}
	}

	return false
}

func AskAndDecryptWithPassphrase(readerEncryted io.Reader) io.Reader {
	scryptIdentity, err := age.NewScryptIdentity(tui.PromptForSecret())
	if err != nil {
		logInternal.Critical.Fatalf("Error while creating passphrase identity: %v\n", err)
	}

	decryptedReader, err := age.Decrypt(readerEncryted, scryptIdentity)
	if err != nil {
		if _, isWrongPassphrase := err.(*age.NoIdentityMatchError); isWrongPassphrase {
			fmt.Printf("Wrong passphrase.\n")
			os.Exit(1)
		} else {
			logInternal.Critical.Fatalf("Error while decrypting key with passphrase: %T\n", err)
		}
	}

	decryptedReaderInBytes, err := io.ReadAll(decryptedReader)
	if err != nil {
		logInternal.Critical.Fatalf("Error while parsing reader: %v\n", err)
	}

	return strings.NewReader(string(decryptedReaderInBytes))
}
