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
		if exists && strings.HasPrefix(href, "/wiki/") {
			title := s.Text()
			url := "https://en.wikipedia.org" + href
			nodes = append(nodes, Node{Title: title, URL: url, Path: []string{title}})
		}
	})
	return nodes, nil
}

func iterativeDeepening(start, target Node, maxDepth int) [][]string {
	foundPaths := make(map[string]bool)
	var results [][]string

	for depth := 0; depth <= maxDepth; depth++ {
		visited := make(map[string]bool)
		var stack []Node
		start.Path = []string{start.Title}
		stack = append(stack, start)

		for len(stack) > 0 {
			current := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if current.Title == target.Title && len(current.Path) <= depth+1 {
				pathStr := strings.Join(current.Path, " -> ")
				if !foundPaths[pathStr] {
					foundPaths[pathStr] = true
					results = append(results, current.Path)
				}
				continue
			}

			if len(current.Path) > depth {
				continue
			}

			if !visited[current.URL] {
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
	}

	return results
}

func main() {
	startNode := Node{Title: "Bandung Institute of Technology", URL: "https://en.wikipedia.org/wiki/Bandung_Institute_of_Technology", Path: []string{"Bandung Institute of Technology"}}
	targetNode := Node{Title: "Indonesia", URL: "https://en.wikipedia.org/wiki/Indonesia", Path: []string{}}

	startTime := time.Now()
	results := iterativeDeepening(startNode, targetNode, 4)
	elapsedTime := time.Since(startTime)

	for _, path := range results {
		fmt.Println("Path found:", strings.Join(path, " -> "))
	}
	fmt.Println("Elapsed time:", elapsedTime)
}
