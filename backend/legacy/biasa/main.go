package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func isInList(str string, list []string) bool {
    for _, v := range list {
        if v == str {
            return true
        }
    }
    return false
}

func FindLink(startURL, transitURL, targetTitle string, depth int, hrefs []string, path []string) bool {
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
	bodyContent.Find("p").Each(func(i int, p *goquery.Selection) {
		p.Find("a").Each(func(j int, s *goquery.Selection) {
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
				if href != "/wiki/Main_Page" && !found && !isInList(href, hrefs) {
					// Continue searching recursively with reduced depth
					hrefs = append(hrefs, href)
					if FindLink(startURL, "https://en.wikipedia.org"+href, targetTitle, depth-1, hrefs, append(path, title)) {
						found = true
						return
					}
				}
			}
		})
	})

	return found
}

func main() {
	startURL := "https://en.wikipedia.org/wiki/Mike_Tyson"
	targetTitle := "Bruce Lee"
	depth := 3
	var hrefs []string
	FindLink(startURL, startURL, targetTitle, depth, hrefs, []string{})
}