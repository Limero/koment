package helper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetPageToJSON[T any](url string, res T) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode < 600 {
		// 4xx or 5xx error
		return fmt.Errorf("got status %q", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(&res)
}

func GetPageBodyString(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode < 600 {
		// 4xx or 5xx error
		return "", fmt.Errorf("got status %q", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}
