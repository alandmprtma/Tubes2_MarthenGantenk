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

func FindLink(startURL, targetTitle string, depth int, hrefs []string, mu *sync.Mutex) {
	queue := []Page{{URL: startURL, Path: []string{}, Depth: depth}}

	for len(queue) > 0 {

		mu.Lock()
		currentPages := make([]Page, 0, 100)
		for i := 0; i < 100 && len(queue) > 0; i++ {
			currentPages = append(currentPages, queue[0])
			queue = queue[1:]
		}
		mu.Unlock()

		var wg sync.WaitGroup
		for _, currentPage := range currentPages {
			wg.Add(1)
			go func(page Page) {
				defer wg.Done()

				if page.Depth == 0 {
					return
				}

				res, err := http.Get(page.URL)
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

				firstPath := doc.Find("title").Text()
				firstPathSplit := strings.Split(firstPath, "-")
				if len(firstPathSplit) > 0 {
					page.Path = append(page.Path, strings.TrimSpace(firstPathSplit[0]))
				}

				content := doc.Find("#mw-content-text")
				content.Find("p").Each(func(i int, p *goquery.Selection) {
					p.Find("a").Each(func(j int, s *goquery.Selection) {
						href, exists := s.Attr("href")
						if exists && strings.HasPrefix(href, "/wiki/") {
							title := s.Text()
							if title == targetTitle {
								result := strings.Join(page.Path, " -> ")
								fmt.Printf("Path: %s -> %s\n", result, title)
								return
							} else if href != "/wiki/Main_Page" && !isInList(href, hrefs) {
								hrefs = append(hrefs, href)
								mu.Lock()
								queue = append(queue, Page{URL: "https://en.wikipedia.org" + href, Path: page.Path, Depth: page.Depth - 1})
								mu.Unlock()
							}
						}
					})
				})
			}(currentPage)
		}
		wg.Wait()
	}
}

func main() {
	// startURL := "https://en.wikipedia.org/wiki/Samsung"
	// targetTitle := "Xiaomi"
	// depth := 2
	startURL := "https://en.wikipedia.org/wiki/Bandung"
	targetTitle := "Japan"
	depth := 3
	var hrefs []string
	var mu sync.Mutex
	start := time.Now()
	FindLink(startURL, targetTitle, depth, hrefs, &mu)
	elapsed := time.Since(start)
	fmt.Printf("Waktu yang dibutuhkan: %s\n", elapsed)
}
