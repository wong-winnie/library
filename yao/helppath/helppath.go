package helppath

import (
	"os"
	"path/filepath"
	"strings"
)

func GetWorkDir() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	return strings.Replace(dir, "\\", "/", -1), nil
}
