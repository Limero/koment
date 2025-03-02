package util

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func GetNumberFromPath(p string) (string, error) {
	re := regexp.MustCompile(`\/(\d+)\/$`)
	matches := re.FindStringSubmatch(p)
	if len(matches) < 2 {
		return "", fmt.Errorf("could not find number in path")
	}
	return matches[1], nil
}

func CachePath(file string) string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(cacheDir, "koment", file)
}
