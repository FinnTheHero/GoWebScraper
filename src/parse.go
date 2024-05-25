package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

/* Replace unwanted strings, <br> & <p> tags */
func ParseChapter(s *goquery.Selection, finalText *string, fileIndex int, unwantedStrings []string) {
	if s.Is("br") || (fileIndex < 78 && s.Is("p")) {
		*finalText += "\n"
	} else if fileIndex > 78 && s.Is("p") {
		text := s.Text()
		if strings.Contains(text, "Chapter") {
			return
		}
		for _, unwantedString := range unwantedStrings {
			text = strings.Replace(text, unwantedString, "", -1)
		}
		if strings.TrimSpace(text) == "" {
			*finalText += "\n"
		} else {
			*finalText += text + "\n"
		}
	} else {
		text := s.Text()
		for _, unwantedString := range unwantedStrings {
			text = strings.Replace(text, unwantedString, "", -1)
		}
		if strings.TrimSpace(text) == "" {
			*finalText += "\n"
		} else {
			*finalText += text
		}
	}
}