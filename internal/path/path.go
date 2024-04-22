package path

import (
	"path/filepath"
	"strings"

	"github.com/lpoirothattermann/storage/internal/constants"
)

func GetNormalizedPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(constants.USER_HOME_DIRECTORY, path[2:])
	}

	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	return path
}
