package utils

import (
	"errors"
	"os"
	"path/filepath"
)

func GetProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "templates")); err == nil {
			return dir, nil
		}
		newDir := filepath.Dir(dir)
		if newDir == dir {
			break
		}
		dir = newDir
	}
	return "", errors.New("корень проекта не найден")
}
