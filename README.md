# WebScraper for TBATE light novel written in GoLang
## About
* Project uses colly framework
* Single threded only, no multithreading as of now
## Extra
* Scraping every chapter takes somewhere between 2 and 3 minutes
* You can further optimize by updating variable `findLastChapter()` function uses to start search from. This will reduce wait time before scrape starts
# Usage
* **Move into `src` directory, inside the project, and run from there**
* **Only use numbers as arguments**

1. **Single Chapter**: Scrape specified chapter
    ```bash
    go run . -single <Chapter>
    ```
2. **Multi Chapter**: Scrape multiple chapters at once 
    * Scrape every chapter
        ```bash
        go run . -multi
        ```
    * Scrape `<From>` till the end
        ```bash
        go run . -multi <From>
        ```
    * Scrape `<From>` to `<to>`
        ```bash
        go run . -multi <From> <To>
        ```
3. **Check** if files already exist
    ```bash
    go run . <Sinlge or Multi> -check
    ```
    > Check can be used on both `Single` and `Multi` mode
