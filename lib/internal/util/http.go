package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:    5,
		IdleConnTimeout: 30 * time.Second,
	},
}

func HTTPGet(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	return httpClient.Do(req)
}

func GetPageToJSON[T any](url string, res T) error {
	resp, err := HTTPGet(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 600 {
		return fmt.Errorf("got status %q with error: %s", resp.Status, body)
	}

	trimmed := bytes.TrimLeft(body, " \t\r\n")
	if len(trimmed) == 0 {
		return fmt.Errorf("empty response body (status %q)", resp.Status)
	}
	if trimmed[0] == '<' {
		return fmt.Errorf("expected JSON but got HTML (status %q) - the API endpoint may be down or blocking requests", resp.Status)
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return fmt.Errorf("%w - response starts with: %s", err, string(trimmed[:min(len(trimmed), 200)]))
	}
	return nil
}

func GetPageBodyString(url string) (string, error) {
	resp, err := HTTPGet(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode < 600 {
		return "", fmt.Errorf("got status %q", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}
