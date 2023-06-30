package helper

import (
	"fmt"
	"strings"
)

func GetFirstBetween(s, from, to string) (string, error) {
	if !strings.Contains(s, from) {
		return "", fmt.Errorf("%q is not in string", from)
	}
	return strings.Split(strings.Split(s, from)[1], to)[0], nil
}

func GetLastBetween(s, from, to string) (string, error) {
	if !strings.Contains(s, from) {
		return "", fmt.Errorf("%q is not in string", from)
	}
	f := strings.Split(s, from)
	return strings.Split(f[len(f)-1], to)[0], nil
}

func CleanText(t string) string {
	t = strings.TrimLeft(t, "\n")
	t = strings.TrimRight(t, "\n")
	t = strings.TrimSpace(t)
	return t
}
