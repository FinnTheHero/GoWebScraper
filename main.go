package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

var help bool
var check bool
var single int
var multi bool
var from = 1
var to = 479

func init() {
	flag.BoolVar(&help, "help", false, "Show help message")
	flag.BoolVar(&check, "check", false, "Check if the files exist to not overwrite")
	flag.IntVar(&single, "single", 0, "Scrape a single chapter instead of all chapters")
	flag.BoolVar(&multi, "multi", false, "Scrape chapters from x to y")
}

func main() {
	flag.Parse()

	if help {
		fmt.Println("Usage: go run main.go [flags]")
		flag.PrintDefaults()
		return
	}

	// Check if the arguments are provided
	args := flag.Args()
	// Check that if arguments are provided they are 2 as needed
	if multi {
		if len(args) > 0 && len(args) < 2 {
			tempFrom, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("Error converting first argument to integer:", err)
				return
			}

			if tempFrom < from || tempFrom > to {
				fmt.Println("The first argument must be greater than ", from, " and less than ", to)
				return
			} else {
				from = tempFrom
			}
		} else if len(args) == 2 {
			tempFrom, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("Error converting first argument to integer:", err)
				return
			}

			if tempFrom < from || tempFrom > to {
				fmt.Println("The first argument must be greater than ", from, " and less than ", to)
				return
			} else {
				from = tempFrom
			}

			tempTo, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Error converting second argument to integer:", err)
				return
			}

			if tempTo > to || tempTo < from {
				fmt.Println("The second argument must be less than ", to, " and greater than ", from)
				return
			} else {
				to = tempTo
			}
		} else if len(args) > 2 {
			fmt.Println("Please provide 2 or less arguments for scraping multiple chapters")
			return
		}
	}

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
		filenames := []string{txtFilename, mdFilename}
		for _, filename := range filenames {
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
			} else {
				file.WriteString(finalText)
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
