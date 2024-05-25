package main

import "time"

/* Add frontmatter to the markdown files */
func Frontmatter(chapterIndex string) string {

	currentDate := time.Now().Format("Jan 2 2006")

	frontmatter := "---\ntitle: 'Chapter " + chapterIndex + "'\nchapter: '" + chapterIndex + "'\ndescription: 'Chapter " + chapterIndex + " of TBATE web-novel'\npubDate: '" + currentDate + "'\nauthor: FinnTheHero\n---"

	return frontmatter
}
