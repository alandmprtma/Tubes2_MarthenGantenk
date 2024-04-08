package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func FindLink(startURL, transitURL, targetTitle string, depth int, path []string) bool {
	if depth < 0 {
		return false
	}

	// Request the HTML page.
	res, err := http.Get(startURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the div with id "bodyContent"
	bodyContent := doc.Find("#bodyContent")
	if bodyContent.Length() == 0 {
		log.Fatal("bodyContent not found")
	}

	// Find the links inside bodyContent
	found := false
	bodyContent.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && strings.HasPrefix(href, "/wiki/") {
			title := s.Text()
			fmt.Printf("Checking: %s\n", title)
			if title == targetTitle {
				fmt.Printf("Found: %s\n", title)
				path = append(path, title)
				fmt.Printf("Path: %v\n", path)
				found = true
				return
			}
			if href != "/wiki/Main_Page" && !found && ("https://en.wikipedia.org"+href) != startURL {
				// Continue searching recursively with reduced depth
				if FindLink(startURL, "https://en.wikipedia.org"+href, targetTitle, depth-1, append(path, title)) {
					found = true
					return
				}
			}
		}
	})

	return found
}

func main() {
	startURL := "https://en.wikipedia.org/wiki/Monster_Hunter"
	targetTitle := "Capcom"
	depth := 3
	FindLink(startURL, startURL, targetTitle, depth, []string{})
}
