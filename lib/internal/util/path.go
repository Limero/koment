package util

import (
	"fmt"
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
