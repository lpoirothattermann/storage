package constants

import (
	"os"
	"path"
)

const VERSION = "0.1"

var userHomeDir, _ = os.UserHomeDir()

var (
	USER_HOME_DIRECTORY                = userHomeDir
	CONFIG_DIRECTORY_PATH              = path.Join(USER_HOME_DIRECTORY, ".config/storage")
	CONFIG_RELATIVE_DEV_DIRECTORY_PATH = "test_data"
)

const PRIVATE_KEY_SIZE_LIMIT = 1 << 24 // 16 MiB

const (
	CONFIG_BASENAME = "storage"
	CONFIG_FILETYPE = "toml"
)

const CONFIG_DEFAULT_LOG_FILE_FILENAME = "storage.log"

const (
	LOG_PREFIX_CRITICAL = "CRITICAL"
	LOG_PREFIX_ERROR    = "ERROR"
	LOG_PREFIX_WARNING  = "WARNING"
	LOG_PREFIX_INFO     = "INFO"
)

// Name of the directory that represents the application in system temporary directory
const APPLICATION_TEMPORARY_DIRECTORY_NAME = "storage"

// Extension that is used to represent backup archives
const ARCHIVE_BACKUP_EXTENSION = ".bk"

// Encrypted, compressed, archived directory with all extensions
const ARCHIVE_EXTENSION = ".tar.gz.age"
