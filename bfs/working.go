package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Page struct {
	URL   string
	Path  []string
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

func FindLink(startURL, targetTitle string, depth int, hrefs []string, wg *sync.WaitGroup) {
	defer wg.Done()

	queue := []Page{{URL: startURL, Path: []string{}, Depth: depth}}
	visited := make(map[string]bool)

	for len(queue) > 0 {
		currentPage := queue[0]
		queue = queue[1:]

		if currentPage.Depth == 0 {
			return
		}

		if visited[currentPage.URL] {
			continue
		}
		visited[currentPage.URL] = true

		res, err := http.Get(currentPage.URL)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		content := doc.Find("#mw-content-text")
		content.Find("p").Each(func(i int, p *goquery.Selection) {
			p.Find("a").Each(func(j int, s *goquery.Selection) {
				href, exists := s.Attr("href")
				if exists && strings.HasPrefix(href, "/wiki/") {
					title := s.Text()
					newPath := append(currentPage.Path, title)
					if title == targetTitle {
						result := strings.Join(newPath, " -> ")
						fmt.Printf("Path: %s\n", result)
						return
					}
					if href != "/wiki/Main_Page" && !isInList(href, hrefs) {
						hrefs = append(hrefs, href)
						wg.Add(1)
						go FindLink("https://en.wikipedia.org"+href, targetTitle, currentPage.Depth-1, hrefs, wg)
					}
				}
			})
		})
	}
}

func main() {
	startURL := "https://en.wikipedia.org/wiki/Samsung"
	targetTitle := "Xiaomi"
	depth := 2
	var hrefs []string
	var wg sync.WaitGroup

	start := time.Now()
	wg.Add(1)
	go FindLink(startURL, targetTitle, depth, hrefs, &wg)

	wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("Waktu yang dibutuhkan: %s\n", elapsed)
}
