package main

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	links := make(map[string]bool)

	url := "https://scrape-me.dreamsofcode.io/"
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	doc, err := html.Parse(response.Body)
	if err != nil {
		panic(err)
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					if strings.HasPrefix(attr.Val, "/") != true {
						fmt.Printf("Skipping %s\n", attr.Val)
						continue
					}

					links[attr.Val] = false
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	for k := range links {

		response, err := http.Get(url + k)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer response.Body.Close()

		doc, err := html.Parse(response.Body)
		if err != nil {
			fmt.Println(err)
			continue
		}

		f(doc)
		links[k] = true
	}

	for k, v := range links {
		fmt.Printf("%s: %t\n", k, v)
	}
}
