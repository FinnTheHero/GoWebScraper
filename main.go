package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector(colly.AllowedDomains("rln.app", "www.rln.app"))

	baseURL := "https://rln.app/the-beginning-after-the-end-535558/chapter-%d"

	baseFileName := "chapter-%d.%s"

	fileIndex := 1

	// Create subfolders for the text files
	mdSubfolder := "md"
	os.MkdirAll(mdSubfolder, os.ModePerm)
	txtSubfolder := "txt"
	os.MkdirAll(txtSubfolder, os.ModePerm)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Scrape the chapter text on each page
	c.OnHTML("div#chapterText", func(e *colly.HTMLElement) {
		unwantedStrings := []string{"SPONSORED CONTENT", "Sponsored Content", "_"}
		finalText := ""

		// Replace unwanted strings and add new lines for <br> tags
		e.DOM.Contents().Each(func(i int, s *goquery.Selection) {
			if s.Is("br") || (fileIndex < 78 && s.Is("p")) {
				finalText += "\n"
			} else if fileIndex > 78 && s.Is("p") {
				text := s.Text()
				for _, unwantedString := range unwantedStrings {
					text = strings.Replace(text, unwantedString, "", -1)
				}
				finalText += text + "\n"
			} else {
				text := s.Text()
				for _, unwantedString := range unwantedStrings {
					text = strings.Replace(text, unwantedString, "", -1)
				}
				finalText += text
			}
		})

		// Create a file name with index
		txtFilename := filepath.Join(txtSubfolder, fmt.Sprintf(baseFileName, fileIndex, "txt"))
		mdFilename := filepath.Join(mdSubfolder, fmt.Sprintf(baseFileName, fileIndex, "md"))

		// Check for the file existence, specify "check" as a first argument
		if len(os.Args) > 1 && os.Args[1] == "check" {
			filenames := []string{txtFilename, mdFilename}
			for _, filename := range filenames {
				if _, err := os.Stat(filename); err == nil {
					fmt.Println("File exists:", filename)
				} else if os.IsNotExist(err) {
					fmt.Println("File doesn't exist:", filename)
				} else {
					fmt.Println("Error finding file:", err)
				}
			}
			return
		}

		// Create a file and write the text to it
		txtFile, txtErr := os.Create(txtFilename)
		mdFile, mdErr := os.Create(mdFilename)

		// Handle errors of creating the file
		if txtErr != nil {
			fmt.Println("Error creating file:", txtErr)
			return
		} else if mdErr != nil {
			fmt.Println("Error creating file:", mdErr)
			return
		} else {
			txtFile.WriteString(finalText)
			txtFile.Close()
			fmt.Println("File saved:", txtFilename)

			mdFile.WriteString(finalText)
			mdFile.Close()
			fmt.Println("File saved:", mdFilename)
		}
	})

	// Loop over for every chapter
	for i := 1; i <= 479; i++ {
		fileIndex = i
		scrapeURL := fmt.Sprintf(baseURL, i)
		c.Visit(scrapeURL)
	}
}
