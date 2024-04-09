package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var wg sync.WaitGroup
var mu sync.Mutex

func isInList(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func FindLink(startURL, transitURL, targetTitle string, depth int, hrefs []string, path []string, ch chan bool) {
	defer wg.Done()

	if depth < 0 {
		ch <- false
		return
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
					mu.Lock()
					path = append(path, title)
					fmt.Printf("Path: %v\n", path)
					mu.Unlock()
					found = true
					return
				}
				if href != "/wiki/Main_Page" && !found && !isInList(href, hrefs) {
					// Continue searching recursively with reduced depth
					hrefs = append(hrefs, href)
					wg.Add(1)
					go FindLink(startURL, "https://en.wikipedia.org"+href, targetTitle, depth-1, hrefs, append(path, title), ch)
				}
			}
		})
	})

	if found {
		ch <- true
	} else {
		ch <- false
	}
}

func main() {
	startURL := "https://en.wikipedia.org/wiki/Mike_Tyson"
	targetTitle := "Bruce Lee"
	depth := 3
	var hrefs []string
	ch := make(chan bool, 1)

	wg.Add(1)
	go FindLink(startURL, startURL, targetTitle, depth, hrefs, []string{}, ch)

	wg.Wait()
	close(ch)

	if <-ch {
		fmt.Println("Target found!")
	} else {
		fmt.Println("Target not found.")
	}
}
