package tui

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/lpoirothattermann/storage/internal/log"
	"golang.org/x/term"
)

/* Below functions directly copied from https://github.com/FiloSottile/age */

func PromptForSecret() string {
	pass, err := readSecret("Enter passphrase:")
	if err != nil {
		log.Critical.Fatalf("could not read passphrase: %v", err)
	}

	return string(pass)
}

// readSecret reads a value from the terminal with no echo. The prompt is ephemeral.
func readSecret(prompt string) (s []byte, err error) {
	err = withTerminal(func(in, out *os.File) error {
		fmt.Fprintf(out, "%s ", prompt)
		defer clearLine(out)
		s, err = term.ReadPassword(int(in.Fd()))
		return err
	})
	return
}

// withTerminal runs f with the terminal input and output files, if available.
// withTerminal does not open a non-terminal stdin, so the caller does not need
// to check stdinInUse.
func withTerminal(f func(in, out *os.File) error) error {
	if runtime.GOOS == "windows" {
		in, err := os.OpenFile("CONIN$", os.O_RDWR, 0)
		if err != nil {
			return err
		}
		defer in.Close()
		out, err := os.OpenFile("CONOUT$", os.O_WRONLY, 0)
		if err != nil {
			return err
		}
		defer out.Close()
		return f(in, out)
	} else if tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0); err == nil {
		defer tty.Close()
		return f(tty, tty)
	} else if term.IsTerminal(int(os.Stdin.Fd())) {
		return f(os.Stdin, os.Stdin)
	} else {
		return fmt.Errorf("standard input is not a terminal, and /dev/tty is not available: %v", err)
	}
}

// clearLine clears the current line on the terminal, or opens a new line if
// terminal escape codes don't work.
func clearLine(out io.Writer) {
	const (
		CUI = "\033["   // Control Sequence Introducer
		CPL = CUI + "F" // Cursor Previous Line
		EL  = CUI + "K" // Erase in Line
	)

	// First, open a new line, which is guaranteed to work everywhere. Then, try
	// to erase the line above with escape codes.
	//
	// (We use CRLF instead of LF to work around an apparent bug in WSL2's
	// handling of CONOUT$. Only when running a Windows binary from WSL2, the
	// cursor would not go back to the start of the line with a simple LF.
	// Honestly, it's impressive CONIN$ and CONOUT$ work at all inside WSL2.)
	fmt.Fprintf(out, "\r\n"+CPL+EL)
}
