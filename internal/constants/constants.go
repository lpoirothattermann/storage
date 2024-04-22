package constants

import "os"

var userHomeDir, _ = os.UserHomeDir()

var USER_HOME_DIRECTORY = userHomeDir + "/"
var CONFIG_DIRECTORY_PATH = USER_HOME_DIRECTORY + ".config/storage/"

const CONFIG_BASENAME = "storage"
const CONFIG_FILETYPE = "toml"

// Name of the directory that represents the application in system temporary directory
const APPLICATION_TEMPORARY_DIRECTORY_NAME = "storage"

// Extension that is used to represent backup archives
const ARCHIVE_BACKUP_EXTENSION = ".bk"

// Encrypted, compressed, archived directory with all extensions
const ARCHIVE_EXTENSION = ".tar.gz.age"
