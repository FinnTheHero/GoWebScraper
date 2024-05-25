package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

var help bool
var check bool
var single int
var multi bool
var from = 1
var to int

func init() {
	flag.BoolVar(&help, "help", false, "Show help message")
	flag.BoolVar(&check, "check", false, "Check if the files exist to not overwrite")
	flag.IntVar(&single, "single", 0, "Scrape a single chapter instead of all chapters")
	flag.BoolVar(&multi, "multi", false, "Scrape chapters from x to y")
}

func main() {
	flag.Parse()

	if help {
		fmt.Println("Usage: go run main.go [flags] [arguments]")
		flag.PrintDefaults()
		return
	}

	FindLastChapter()

	// Check if the arguments are provided
	args := flag.Args()
	if multi {
		HandleMulti(args, &from, &to)
	}

	if single != 0 {
		HandleSingle(single, from, to)
	}

	c := colly.NewCollector(colly.AllowedDomains("rln.app", "www.rln.app"))

	baseURL := "https://rln.app/the-beginning-after-the-end-535558/chapter-%d"

	baseFileName := "chapter-%d.%s"

	fileIndex := 1

	// Create subfolders for the text files
	mdSubfolder := "../md"
	os.MkdirAll(mdSubfolder, os.ModePerm)
	txtSubfolder := "../txt"
	os.MkdirAll(txtSubfolder, os.ModePerm)

	c.OnRequest(func(r *colly.Request) {
		fmt.Print(r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println(" | ", r.StatusCode)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Scrape the chapter text on each page
	c.OnHTML("div#chapterText", func(e *colly.HTMLElement) {
		unwantedStrings := []string{"SPONSORED CONTENT", "Sponsored Content", "_", "                            "}
		// Add frontmatter first then rest of the text
		finalText := ""

		e.DOM.Contents().Each(func(i int, s *goquery.Selection) {
			ParseChapter(s, &finalText, fileIndex, unwantedStrings)
		})

		// Create a file name with index
		txtFilename := filepath.Join(txtSubfolder, fmt.Sprintf(baseFileName, fileIndex, "txt"))
		mdFilename := filepath.Join(mdSubfolder, fmt.Sprintf(baseFileName, fileIndex, "md"))

		// Check for the file existence
		filenames := []string{txtFilename, mdFilename}
		for i, filename := range filenames {
			if check {
				if _, err := os.Stat(filename); err == nil {
					fmt.Println("File exists:", filename)
					continue
				} else if os.IsNotExist(err) {
					fmt.Println("File doesn't exist:", filename)
				} else {
					fmt.Println("Error finding file:", err)
					return
				}
			}

			file, err := os.Create(filename)

			if err != nil {
				fmt.Println("Error creating file:", err)
				return
			}

			if i == 0 {
				file.WriteString(finalText)
				file.Close()
				fmt.Println("File saved:", filename)
			} else if i == 1 {
				file.WriteString(Frontmatter(strconv.Itoa(fileIndex)) + finalText)
				file.Close()
				fmt.Println("File saved:", filename)
			}

		}
	})

	if single != 0 {
		// Scrape a single chapter
		fmt.Println("Scraping chapter:", single)
		fileIndex = single
		scrapeURL := fmt.Sprintf(baseURL, single)
		c.Visit(scrapeURL)
	} else if single == 0 && multi {
		// Loop over for every chapter
		fmt.Println("Scraping chapters from ", from, " to ", to)
		for i := from; i <= to; i++ {
			fileIndex = i
			scrapeURL := fmt.Sprintf(baseURL, i)
			c.Visit(scrapeURL)
		}
	}
}
