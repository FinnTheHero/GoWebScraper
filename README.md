# WebScraper for TBATE light novel written in Go
## About
* Project uses colly framework

## Extra
* Scraping every chapter takes somewhere between 2 and 3 minutes
* You can further optimize by updating variable `findLastChapter()` function uses to start search from. This will reduce wait time before scrape starts
# Usage
1. **Single Chapter**: Scrape specified chapter
    ```bash
    go run main.go -single <Chapter To Scrape>
    ```
2. **Multi Chapter**: Scrape multiple chapters at once 
    * Scrape every chapter
    ```bash
    go run main.go -multi
    ```
    * Scrape `<From>` till the end
    ```bash
    go run main.go -multi <From>
    ```
    * Scrape `<From>` `<to>`
    ```bash
    go run main.go -multi <From> <To>
    ```
3. **Check** if files already exist
    ```bash
    go run main.go -multi -check
    ```
