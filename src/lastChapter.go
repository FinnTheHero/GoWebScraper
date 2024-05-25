package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func FindLastChapter() {
	pageIndex := 482

	exitLoop := false

	fmt.Println("Searching for the last chapter...")
	fmt.Println("---------------------------------")

	c := colly.NewCollector(colly.AllowedDomains("rln.app", "www.rln.app"))

	c.OnError(func(r *colly.Response, err error) {
		fmt.Print(r.Request.URL)
		fmt.Println(" | ", r.StatusCode)
		fmt.Println("Error: ", err)

		if r.StatusCode != 200 {
			fmt.Println("Last chapter: ", pageIndex-1)
			fmt.Println("---------------------------------")
			to = pageIndex - 1
			exitLoop = true
			return
		}
	})

	for {
		if exitLoop {
			break
		}

		c.Visit(fmt.Sprintf("https://rln.app/the-beginning-after-the-end-535558/chapter-%d", pageIndex))

		pageIndex++
	}
}
