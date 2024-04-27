package log

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/lpoirothattermann/storage/internal/constants"
)

var (
	logFile  *os.File
	Info     *log.Logger
	Warning  *log.Logger
	Error    *log.Logger
	Critical *log.Logger
)

func Initialization(logFilePath string) error {
	var err error
	logFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	Critical = log.New(multiWriter, constants.LOG_PREFIX_CRITICAL+": ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(multiWriter, constants.LOG_PREFIX_ERROR+": ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(multiWriter, constants.LOG_PREFIX_WARNING+": ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(multiWriter, constants.LOG_PREFIX_INFO+": ", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}

func Close() error {
	if logFile == nil {
		return errors.New("Can't be closed, has not been initialized")
	}

	if err := logFile.Close(); err != nil {
		return err
	}

	logFile = nil
	Info = nil
	Warning = nil
	Error = nil
	Critical = nil

	return nil
}
