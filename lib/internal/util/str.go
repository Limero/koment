package util

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetLastBetween(s, from, to string) (string, error) {
	if !strings.Contains(s, from) {
		return "", fmt.Errorf("%q is not in string", from)
	}
	f := strings.Split(s, from)
	return strings.Split(f[len(f)-1], to)[0], nil
}

func CleanHTML(t string) string {
	if doc, err := goquery.NewDocumentFromReader(strings.NewReader(t)); err == nil {
		doc.Find("p").Each(func(_ int, s *goquery.Selection) {
			s.AfterHtml("\n")
		})
		return doc.Text()
	}
	return t
}
