package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Node struct {
	Title string
	URL   string
	Path  []string
}

var httpClient = &http.Client{
	Timeout: time.Second * 10, // Set timeout
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 10,
	},
}

func fetchLinks(url string) ([]Node, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var nodes []Node
	doc.Find("#bodyContent p a").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && strings.HasPrefix(href, "/wiki/") && !strings.Contains(href, ":") {
			title := s.Text()
			fullURL := "https://en.wikipedia.org" + href
			nodes = append(nodes, Node{Title: title, URL: fullURL, Path: []string{title}})
		}
	})
	return nodes, nil
}

func iterativeDeepening(start, target Node, maxDepth int) [][]string {
	var results [][]string
	var foundFirstDepth bool

	for depth := 0; depth <= maxDepth && !foundFirstDepth; depth++ {
		visited := make(map[string]bool)
		var stack []Node
		start.Path = []string{start.Title}
		stack = append(stack, start)

		for len(stack) > 0 {
			current := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			fmt.Println("Checking path:", strings.Join(current.Path, " -> "))

			if current.Title == target.Title && len(current.Path) <= depth+1 {
				results = append(results, current.Path)
				foundFirstDepth = true
			}

			if len(current.Path) > depth || visited[current.URL] {
				continue
			}

			visited[current.URL] = true
			neighbors, err := fetchLinks(current.URL)
			if err != nil {
				log.Println(err)
				continue
			}

			for _, neighbor := range neighbors {
				if !visited[neighbor.URL] {
					neighbor.Path = append([]string(nil), current.Path...)
					neighbor.Path = append(neighbor.Path, neighbor.Title)
					stack = append(stack, neighbor)
				}
			}
		}
	}

	return results
}

func main() {
	startNode := Node{Title: "Joko Widodo", URL: "https://en.wikipedia.org/wiki/Joko_Widodo", Path: []string{"Joko Widodo"}}
	targetNode := Node{Title: "Bandung Institute of Technology", URL: "https://en.wikipedia.org/wiki/Bandung_Institute_of_Technology", Path: []string{}}

	startTime := time.Now()
	results := iterativeDeepening(startNode, targetNode, 6)
	elapsedTime := time.Since(startTime)
	fmt.Println("\nPath(s) found:")
	pathCounter := 1
	for _, path := range results {
		fmt.Printf("%d. %s\n", pathCounter, strings.Join(path, " -> "))
		pathCounter++
	}
	fmt.Println("Elapsed time:", elapsedTime)
}
