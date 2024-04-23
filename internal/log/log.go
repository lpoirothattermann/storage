package log

import (
	"errors"
	"io"
	"log"
	"os"
)

var (
	logFile  *os.File
	Info     *log.Logger
	Warning  *log.Logger
	Error    *log.Logger
	Critical *log.Logger
)

func Initialization(logFilePath string) error {
	LogFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}

	multiWriter := io.MultiWriter(os.Stdout, LogFile)

	Info = log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(multiWriter, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(multiWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	Critical = log.New(multiWriter, "CRITICAL: ", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}

func Close() error {
	if logFile != nil {
		errors.New("Can't be closed, has not been initialized")
	}

	if err := logFile.Close(); err != nil {
		return err
	}

	return nil
}
