package main

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func scrapeURL(url string, list map[string]bool) {
	fmt.Println("Scraping", url)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()

	doc, err := html.Parse(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	findAnchors(doc, list)
}

func findAnchors(n *html.Node, list map[string]bool) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" && strings.HasPrefix(attr.Val, "/") {
				list[strings.TrimPrefix(attr.Val, "/")] = false
				continue
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findAnchors(c, list)
	}
}

func main() {
	links := make(map[string]bool)

	scrapeURL("https://scrape-me.dreamsofcode.io/", links)
	for link := range links {
		scrapeURL("https://scrape-me.dreamsofcode.io/"+link, links)
	}

	for link := range links {
		fmt.Println(link)
	}

	// Send a link
	// scrape the link
	// send sub links
	// repeat
}
