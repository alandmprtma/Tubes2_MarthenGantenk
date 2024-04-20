package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Node struct {
	Title string
	URL   string
	Path  []string
}

var httpClient = &http.Client{
	Timeout: time.Second * 10,
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 100,
	},
}

func fetchLinks(url string, ch chan<- []Node, errCh chan<- error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Connection", "keep-alive")

	res, err := httpClient.Do(req)
	if err != nil {
		errCh <- err
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		errCh <- fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		errCh <- err
		return
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
	ch <- nodes
}

func worker(urlQueue <-chan string, resultQueue chan<- []Node, errQueue chan<- error, wg *sync.WaitGroup) {
	for url := range urlQueue {
		ch := make(chan []Node)
		errCh := make(chan error)
		go fetchLinks(url, ch, errCh)
		select {
		case result := <-ch:
			resultQueue <- result
		case err := <-errCh:
			errQueue <- err
		}
	}
	wg.Done()
}

func iterativeDeepening(start, target Node, maxDepth int) [][]string {
	var results [][]string
	pathSet := make(map[string]bool) // Map to track unique paths
	urlQueue := make(chan string, 10)
	resultQueue := make(chan []Node, 10)
	errQueue := make(chan error, 10)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker(urlQueue, resultQueue, errQueue, &wg)
	}

	var foundDepth int = maxDepth + 1
	for depth := 0; depth <= maxDepth && depth <= foundDepth; depth++ {
		visited := make(map[string]bool)
		var stack []Node
		start.Path = []string{start.Title}
		stack = append(stack, start)

		for len(stack) > 0 {
			current := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			pathKey := strings.Join(current.Path, " -> ")
			if current.Title == target.Title && !pathSet[pathKey] {
				fmt.Printf("Path found: %s\n", pathKey) // Print the path immediately when found
				results = append(results, current.Path)
				pathSet[pathKey] = true
				// if depth < foundDepth {
				// 	foundDepth = depth
				// }
			}

			if len(current.Path) > depth || visited[current.URL] {
				continue
			}

			visited[current.URL] = true
			urlQueue <- current.URL

			select {
			case neighbors := <-resultQueue:
				for _, neighbor := range neighbors {
					if !visited[neighbor.URL] {
						neighbor.Path = append([]string(nil), current.Path...)
						neighbor.Path = append(neighbor.Path, neighbor.Title)
						stack = append(stack, neighbor)
					}
				}
			case err := <-errQueue:
				log.Println(err)
			}
		}
	}

	close(urlQueue)
	wg.Wait()
	close(resultQueue)
	close(errQueue)
	return results
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	startNode := Node{Title: "Bandung Institute of Technology", URL: "https://en.wikipedia.org/wiki/Bandung_Institute_of_Technology", Path: []string{"Bandung Institute of Technology"}}
	targetNode := Node{Title: "Japan", URL: "https://en.wikipedia.org/wiki/Japan", Path: []string{}}

	fmt.Println("Starting search...")
	startTime := time.Now()
	results := iterativeDeepening(startNode, targetNode, 6)
	elapsedTime := time.Since(startTime)

	fmt.Println("\nFinal unique paths found:")
	pathCounter := 1
	for _, path := range results {
		fmt.Printf("%d. %s\n", pathCounter, strings.Join(path, " -> "))
		pathCounter++
	}

	fmt.Println("Search completed.")
	fmt.Println("Elapsed time:", elapsedTime)
}
