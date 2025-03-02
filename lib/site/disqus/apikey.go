package disqus

import (
	"github.com/limero/koment/lib/internal/util"
)

// Return API key if it exists, otherwise fetch it from disqus
func (s Disqus) getApiKey() (string, error) {
	apiKeyFile := util.CachePath("disqus-api-key.txt")
	apiKey, err := util.ReadFileIfExists(apiKeyFile)
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

	return apiKey, util.WriteFile(apiKey, apiKeyFile)
}

func (s Disqus) fetchApiKey() (string, error) {
	body, err := util.GetPageBodyString("https://disqus.disqus.com/embed.js")
	if err != nil {
		return "", err
	}
	disqusVersion, err := util.GetLastBetween(body, "lounge.load.", ".js")
	if err != nil {
		return "", err
	}

	url := "https://c.disquscdn.com/next/embed/lounge.load." + disqusVersion + ".js"
	body, err = util.GetPageBodyString(url)
	if err != nil {
		return "", err
	}
	bundleVersion, err := util.GetLastBetween(body, "common.bundle.", ".js")
	if err != nil {
		return "", err
	}

	url = "https://c.disquscdn.com/next/embed/common.bundle." + bundleVersion + ".js"
	body, err = util.GetPageBodyString(url)
	if err != nil {
		return "", err
	}
	apiKey, err := util.GetLastBetween(body, "embedAPI:\"", "\"")
	if err != nil {
		return "", err
	}
	return apiKey, nil
}
