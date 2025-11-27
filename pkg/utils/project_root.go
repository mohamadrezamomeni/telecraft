package utils

import (
	"os"
	"path/filepath"
)

func GetRootOfProject() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	var previous any
	for previous != dir {
		if _, err = os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		previous = dir
		dir = filepath.Dir(dir)
	}

	return "", os.ErrNotExist
}
