package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-shiori/go-epub"
)

/* Add frontmatter to the markdown files */
func Frontmatter(chapterIndex string) string {

	currentDate := time.Now().Format("Jan 2 2006")

	frontmatter := "---\ntitle: 'Chapter " + chapterIndex + "'\nchapter: '" + chapterIndex + "'\ndescription: 'Chapter " + chapterIndex + " of TBATE web-novel'\npubDate: '" + currentDate + "'\nauthor: FinnTheHero\n---"

	return frontmatter
}

func InsertString(original string, insert string, index int) string {
	if index < 0 || index > len(original) {
		return original
	}
	return original[:index] + insert + original[index:]
}

/* Create EPUB sections */
func CreateEPUBSection(e *epub.Epub, text string, cssPath string, chapterIndex int) {
	sectionTitle := `<h2 class="chapter-title">Chapter ` + strconv.Itoa(chapterIndex) + `</h2>`
	sectionBody := ""

	paragraphs := strings.Split(text, "\n")

	for _, paragraph := range paragraphs {
		t := strings.ReplaceAll(paragraph, "	", "")
		t = strings.ReplaceAll(t, "\n", "")
		sectionBody += `<p>	   ` + t + `</p>`
	}

	// sectionBody = `<p class="dialog-color">` + strings.ReplaceAll(text, "\n", "<br/>") + `</p>`

	finalText := sectionTitle + sectionBody

	_, err := e.AddSection(finalText, "Chapter "+strconv.Itoa(chapterIndex), "Chapter "+strconv.Itoa(chapterIndex), cssPath)
	if err != nil {
		log.Fatalf("Error adding section: %v", err)
	}
}

/* Create EPUB */
func CreateEPUB() {

	epubPath := "../epub/tbate.epub"

	e, err := epub.NewEpub("The Beginning After The End")
	if err != nil {
		log.Fatalf("Error creating EPUB: %v", err)
	}

	e.SetAuthor("TurtleMe")

	e.SetTitle("The Beginning After The End")

	e.SetDescription("King Grey has unsurpassed strength, wealth and authority in a world where military abilities rule. However, loneliness remains with those who have great power. Under the glamorous appearance of a powerful king, there is a shell of a man, devoid of purpose and will. Reincarnated in a new world filled with magic and monsters, the king gets a second chance to live his life anew. However, correcting the mistakes of the past will not be his only task. Under the peace and prosperity of the new world, there is a hidden current that threatens to destroy everything for which he worked, questioning his role and the reason for being born again.")

	// Add CSS
	cssPath, err := e.AddCSS("../epub/src/style.css", "style.css")
	if err != nil {
		log.Fatalf("Error adding CSS: %v", err)
	}

	// Add Arial font
	e.AddFont("../epub/src/arial.ttf", "arial.ttf")

	// Add cover image
	imgPath, err := e.AddImage("../epub/src/cover.jpg", "cover")
	if err != nil {
		log.Fatalf("Error adding image: %v", err)
	}
	e.SetCover(imgPath, cssPath)

	// Get the amount of files in the /txt directory
	files, err := os.ReadDir("../txt")
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	for i := 0; i < len(files); i++ {
		// Skip non-txt files
		if strings.Split(files[i].Name(), ".")[1] != "txt" {
			log.Printf("Not a txt file: %v", files[i].Name())
			return
		}

		path := "../txt/chapter-" + strconv.Itoa(i+1) + ".txt"

		// Read the file
		textBytes, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("Error reading file: %v", err)
		}

		// Convert byte array into string
		text := string(textBytes)

		// Create a new EPUB section
		CreateEPUBSection(e, text, cssPath, i+1)
	}

	e.Write(epubPath)
	fmt.Print("EPUB created successfully: ", epubPath)
}

/* Replace unwanted strings, <br> & <p> tags */
func ParseChapter(s *goquery.Selection, parsedText *string, fileIndex int, unwantedStrings []string) {
	if s.Is("br") || (fileIndex < 78 && s.Is("p")) {
		*parsedText += "\n"
	} else if fileIndex > 78 && s.Is("p") {
		text := s.Text()
		if strings.Contains(text, "Chapter") {
			return
		}
		for _, unwantedString := range unwantedStrings {
			text = strings.Replace(text, unwantedString, "", -1)
		}
		if strings.TrimSpace(text) == "" {
			*parsedText += "\n"
		} else {
			*parsedText += text + "\n"
		}
	} else {
		text := s.Text()
		for _, unwantedString := range unwantedStrings {
			text = strings.Replace(text, unwantedString, "", -1)
		}
		if strings.TrimSpace(text) == "" {
			*parsedText += "\n"
		} else {
			*parsedText += text
		}
	}
}
