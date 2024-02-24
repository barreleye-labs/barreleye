package common

import (
	"os"
	"path/filepath"
)

func GetProjectRootPath() string {
	f, _ := os.Getwd()
	return filepath.Join(filepath.Dir(f), filepath.Base(f))
}
