package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Page struct {
	URL  string
	Path []string
}

func isInList(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func FindLink(startURL, targetTitle string, depth int) []string {
	queue := []Page{{URL: startURL, Path: []string{}}}
	visited := make(map[string]bool)
	var result []string

	for len(queue) > 0 {
		currentPage := queue[0]
		queue = queue[1:]

		// Skip if the page is already visited
		if visited[currentPage.URL] {
			continue
		}

		// Request the HTML page
		res, err := http.Get(currentPage.URL)
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

		// Find the links inside the page
		bodyContent := doc.Find("#bodyContent")
		bodyContent.Find("p").Each(func(i int, p *goquery.Selection) {
			p.Find("a").Each(func(j int, s *goquery.Selection) {
				href, exists := s.Attr("href")
				if exists && strings.HasPrefix(href, "/wiki/") {
					title := s.Text()
					newPath := append(currentPage.Path, title)
					if title == targetTitle {
						result = newPath
						return
					}
					if href != "/wiki/Main_Page" && !isInList(href, visited) {
						queue = append(queue, Page{URL: "https://en.wikipedia.org" + href, Path: newPath})
					}
				}
			})
		})

		visited[currentPage.URL] = true
	}

	return result
}

func main() {
	startURL := "https://en.wikipedia.org/wiki/Samsung"
	targetTitle := "Xiaomi"
	depth := 3
	path := FindLink(startURL, targetTitle, depth)
	if len(path) > 0 {
		fmt.Printf("Path found: %v\n", path)
	} else {
		fmt.Println("Path not found")
	}
}
