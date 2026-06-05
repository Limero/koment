package util

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func CleanHTML(t string) string {
	if doc, err := goquery.NewDocumentFromReader(strings.NewReader(t)); err == nil {
		doc.Find("p").Each(func(_ int, s *goquery.Selection) {
			s.AfterHtml("\n")
		})
		return doc.Text()
	}
	return t
}
