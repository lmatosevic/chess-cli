package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func HomeFilePath(path string) string {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		userHomeDir = "."
	}

	sep := "/"
	if strings.HasPrefix(path, "/") || strings.HasSuffix(userHomeDir, "/") {
		sep = ""
	}

	return fmt.Sprintf("%s%s%s", userHomeDir, sep, path)
}

func ReadFromFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func WriteToFile(path string, data []byte) error {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.Write(data); err != nil {
		return err
	}

	return nil
}

func DeleteFile(path string) error {
	err := os.Remove(path)
	return err
}

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
