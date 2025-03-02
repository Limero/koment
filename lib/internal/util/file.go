package util

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func ReadFileIfExists(path string) (string, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", nil
	} else if err != nil {
		return "", err
	}

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var content string
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return strings.TrimSpace(content), nil
}

func WriteFile(content, destination string) error {
	/*
		Create file at destination with content. If file exists, it will be overwritten.
	*/
	parentDir := filepath.Dir(destination)
	err := os.MkdirAll(parentDir, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return err
	}

	return nil
}
