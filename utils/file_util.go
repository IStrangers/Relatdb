package utils

import (
	"os"
	"strings"
)

func ConcatFilePaths(paths ...string) string {
	return strings.Join(paths, string(os.PathSeparator))
}
