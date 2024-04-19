package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Page struct {
	URL  string
	Path []string
	Depth int
}

func isInList(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func FindLink(startURL, targetTitle string, depth int, hrefs []string) {
	queue := []Page{{URL: startURL, Path: []string{}, Depth: depth}}

	for len(queue) > 0 {
		
		currentPage := queue[0]
		queue = queue[1:]

		if (currentPage.Depth == 0) {
			return
		}
		// fmt.Println(currentPage.URL)

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

		firstPath := doc.Find("title").Text()
		firstPathSplit := strings.Split(firstPath, "-")
		if len(firstPathSplit) > 0 {
			currentPage.Path = append(currentPage.Path, strings.TrimSpace(firstPathSplit[0]))
		} 
		
		// Find the links inside the page
		content := doc.Find("#mw-content-text")
		content.Find("p").Each(func(i int, p *goquery.Selection) {
			p.Find("a").Each(func(j int, s *goquery.Selection) {
				href, exists := s.Attr("href")
				// fmt.Println(href)
				if exists && strings.HasPrefix(href, "/wiki/") {
					title := s.Text()
					// fmt.Printf("Checking: %s\n", title)
					// newPath := append(currentPage.Path, title)
					if title == targetTitle {
						result := strings.Join(currentPage.Path, " -> ")
						fmt.Printf("Path: %s -> %s\n", result, title)
						return
					} else if href != "/wiki/Main_Page" && !isInList(href, hrefs) {
						hrefs = append(hrefs, href)
						queue = append(queue, Page{URL: "https://en.wikipedia.org" + href, Path: currentPage.Path, Depth: currentPage.Depth-1})
					}
				}
			})
		})
	}
}

func main() {
	startURL := "https://en.wikipedia.org/wiki/Basketball"
	targetTitle := "Bandung Institute of Technology"
	depth := 3
	var hrefs []string
	start := time.Now()
	FindLink(startURL, targetTitle, depth, hrefs)
	elapsed := time.Since(start)
	fmt.Printf("Waktu yang dibutuhkan: %s\n", elapsed)
}
