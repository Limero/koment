package disqus

import "github.com/limero/koment/lib/internal/helper"

// Return API key if it exists, otherwise fetch it from disqus
func (s Disqus) getApiKey() (string, error) {
	apiKeyFile := helper.CachePath("disqus-api-key.txt")
	apiKey, err := helper.ReadFileIfExists(apiKeyFile)
	if err != nil {
		return "", err
	}
	if apiKey != "" {
		return apiKey, nil
	}

	apiKey, err = s.fetchApiKey()
	if err != nil {
		return "", err
	}

	return apiKey, helper.WriteFile(apiKey, apiKeyFile)
}

func (s Disqus) fetchApiKey() (string, error) {
	body, err := helper.GetPageBodyString("https://disqus.disqus.com/embed.js")
	if err != nil {
		return "", err
	}
	disqusVersion, err := helper.GetLastBetween(body, "lounge.load.", ".js")
	if err != nil {
		return "", err
	}

	url := "https://c.disquscdn.com/next/embed/lounge.load." + disqusVersion + ".js"
	body, err = helper.GetPageBodyString(url)
	if err != nil {
		return "", err
	}
	bundleVersion, err := helper.GetLastBetween(body, "common.bundle.", ".js")
	if err != nil {
		return "", err
	}

	url = "https://c.disquscdn.com/next/embed/common.bundle." + bundleVersion + ".js"
	body, err = helper.GetPageBodyString(url)
	if err != nil {
		return "", err
	}
	apiKey, err := helper.GetLastBetween(body, "embedAPI:\"", "\"")
	if err != nil {
		return "", err
	}
	return apiKey, nil
}
