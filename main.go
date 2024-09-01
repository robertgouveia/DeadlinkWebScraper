package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

func scrape(baseURL string, url string, deadLinks []string, links map[string]bool, mux *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	// this function will initiate the scraping of a url by parsing the body of the url

	fullURL := baseURL + url // merges the base url and the url to be scraped
	//fmt.Println("Scraping", fullURL)

	response, err := http.Get(fullURL) // sends a get request to the url (brings back html.response)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()
	// closes the response body when the function ends
	// this will prevent the body from being leaked

	if response.StatusCode != 200 {
		// if the url does not exist
		fmt.Println("Dead link", fullURL)
		deadLinks = append(deadLinks, fullURL)
	}

	doc, err := html.Parse(response.Body)
	// parses the html response body, this returns a node list
	if err != nil {
		fmt.Println(err)
		return
	}

	// scrape the url and find all the links
	newLinks := findAnchors(doc, baseURL, links, mux) // finds all the links in the html node

	// each new sub link we now need to scrape them
	for _, link := range newLinks {
		mux.Lock() // lock the map
		// if the link has not been scraped
		if !links[link] {
			links[link] = true
			wg.Add(1) // add one to the wait group
			go scrape(baseURL, link, deadLinks, links, mux, wg)
			// recursively scrape the link
		}
		mux.Unlock() // unlock the map
	}
}

func findAnchors(n *html.Node, baseURL string, links map[string]bool, mux *sync.Mutex) []string {
	var newLinks []string // this will hold and return any new links found

	// this function will find all the links in the html node

	if n.Type == html.ElementNode && n.Data == "a" {
		// if the node is an element and a link
		for _, attr := range n.Attr {
			// loop through attributes (class / id / href)
			if attr.Key == "href" && (strings.HasPrefix(attr.Val, "/")) {
				// we only want the href and a relative link
				// if the url does hold the base url
				relativePath := strings.TrimPrefix(attr.Val, baseURL)
				// remove the base url from the link

				//check if it has already been scraped
				mux.Lock()
				// ensure we are the only goroutine accessing the map
				if _, exists := links[relativePath]; !exists {
					// getting the bool value of the link, if false
					links[relativePath] = false
					// add the link to the map
					newLinks = append(newLinks, relativePath)
				}
				mux.Unlock()
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		//this goes from one node to the next
		newLinks = append(newLinks, findAnchors(c, baseURL, links, mux)...)
		// recursively sends in nodes to check for links and then appends them to the newLinks slice
	}

	return newLinks
	// returns the new links found
}

func main() {
	var deadLinks []string

	baseURL := "https://scrape-me.dreamsofcode.io/"
	// base url is the url to be scraped
	links := make(map[string]bool)
	// links is a map with string key and bool value

	var wg sync.WaitGroup
	// wait group is a thread safe counter
	// when the counter is 0, the wait group is done
	// keeps track of the number of goroutines

	var mux sync.Mutex
	// mutex locks the map
	// no two goroutines can access the map at the same time
	// this is important because we are using a map in multiple goroutines

	// deadlocks can occur if two goroutines try to access the map at the same time
	// deadlocks are when two goroutines are waiting for each other to finish

	initialPath := "/"
	// this is to be scraped first

	wg.Add(1) // add one to the wait group
	// this means that there is one goroutine running

	// scrape the initial base url with a / path
	go scrape(baseURL, initialPath, deadLinks, links, &mux, &wg)

	wg.Wait() // wait for all goroutines to finish

	for link := range links {
		//print all links that have been scraped

		if links[link] {
			fmt.Println(link)
		}
	}

	for deadLinks := range deadLinks {
		// print all links that are dead
		fmt.Println(deadLinks)
	}
}
