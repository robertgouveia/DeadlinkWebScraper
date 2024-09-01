
# Dead links Web Scraper

This scrapes HTML based sites for links, sub links. This will return any links that do not respond with a 200 OK response.
## Installation

Install golang first

```
https://go.dev/dl/
```

Clone the project and run main

```bash
  git clone https://github.com/robertgouveia/DeadlinkWebScraper
  cd DeadlinkWebScraper
  go run main.go
```

    
## Response Example

In the code on line 101 there is a url

```go
 baseURL := "https://scrape-me.dreamsofcode.io/" 
```

This will be the url to scrape.

The response for running this may be different each time in order of dead links due to concurrency however the dead links of this page are listed like so:

```bash
Dead link /nevermind
Dead link /in-utero
Dead link /venus
Dead link /earth
Dead link /mars
Dead link /recursion
Dead link /mountain
Dead link /teapot
Dead link /busted
/nevermind
/in-utero
/venus
/earth
/mars
/recursion
/mountain
/teapot
/busted
```
## Roadmap

- Support for Javascript based pages

- Using packages that can assist with web scraping

- Optimise the use of concurrency

