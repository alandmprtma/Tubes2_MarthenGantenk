package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// define structure to represent node in graph
type Node struct {
	Title string
	URL   string
	Path  []string
}

type Result struct {
	Results           []string `json:"results"`
	ArticlesChecked   int      `json:"articlesChecked"`
	ArticlesTraversed int      `json:"articlesTraversed"`
	NumberPath        int      `json:"numberPath"`
	ElapsedTime       float64  `json:"elapsedTime"`
}

var (
	// create an https client with 10 seconds timeout and a transport with a specific max idle connections
	httpClient = &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 100,
		},
	}
	urlCache = make(map[string][]Node)  // cache for storing request result
	sem      = make(chan struct{}, 500) // semaphore to limit the number of http request
)

// function to fetch links from a given url
func fetchLinks(urlString string, ch chan<- []Node, errCh chan<- error) {
	// reserve a spot in the semaphore to limit concurrent HTTP requests
	sem <- struct{}{}
	defer func() { <-sem }() // release the semaphore when it's done

	// cek if the result url has already in cache
	if nodes, ok := urlCache[urlString]; ok {
		ch <- nodes
		return
	}

	// create a new http get request
	req, _ := http.NewRequest("GET", urlString, nil)
	req.Header.Set("Connection", "keep-alive")

	// perform the http request
	res, err := httpClient.Do(req)
	if err != nil {
		errCh <- err
		return
	}
	defer res.Body.Close() // ensure the response body is closed after processing

	// check if the response status code indicates success
	if res.StatusCode != 200 {
		errCh <- fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
		return
	}

	// parse the html content to extract links
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		errCh <- err
		return
	}

	// store extracted nodes
	var nodes []Node
	doc.Find("#mw-content-text a").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && strings.HasPrefix(href, "/wiki/") && !strings.Contains(href, ":") {
			// Extract the title from the URL
			title := strings.TrimPrefix(href, "/wiki/")
			title, err = url.QueryUnescape(title) // Decode URL-encoded characters
			if err != nil {
				errCh <- err
				return
			}
			title = strings.ReplaceAll(title, "_", " ") // Replace underscores with spaces
			fullURL := "https://en.wikipedia.org" + href
			nodes = append(nodes, Node{Title: title, URL: fullURL, Path: []string{title}}) // create a new node with the link's details
		}
	})

	// store the result to cache
	urlCache[urlString] = nodes
	ch <- nodes
}

// worker function to process urls from a queue
func worker(urlQueue <-chan string, resultQueue chan<- []Node, errQueue chan<- error, wg *sync.WaitGroup) {
	for url := range urlQueue {

		// create channel for communication
		ch := make(chan []Node)
		errCh := make(chan error)

		// fetch link from the given url
		go fetchLinks(url, ch, errCh)

		// process the result or error based on the response
		select {
		case result := <-ch:
			resultQueue <- result
		case err := <-errCh:
			errQueue <- err
		}
	}

	// signal that this worker is done
	wg.Done()
}

// iterative deepening search to find unique paths from a start node to a target node
func iterativeDeepeningAll(start, target Node, maxDepth int) ([][]string, int, int, int, float64) {
	startTime := time.Now()
	var results [][]string                // to store the final result
	pathSet := make(map[string]bool)      // to ensure unique paths
	urlQueue := make(chan string, 500)    // channel for the urls to be processed
	resultQueue := make(chan []Node, 500) // channel for the results
	errQueue := make(chan error, 500)     // channel for the errors

	// create multiple worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < 30; i++ { // set 30 workers
		wg.Add(1)
		go worker(urlQueue, resultQueue, errQueue, &wg) // start a new worker
	}

	var foundDepth int = maxDepth + 1 // for tracking the first depth the path is found

	// counters for articles checked and traversed
	var articlesChecked int
	var articlesTraversed int

	// iterate through increasing depth
	for depth := 0; depth <= maxDepth && depth <= foundDepth; depth++ {
		visited := make(map[string]bool) // track visited nodes
		var stack []Node

		// set the initial path for the start node
		start.Path = []string{start.Title}
		stack = append(stack, start)

		for len(stack) > 0 {
			current := stack[len(stack)-1] // get the top of the stack
			stack = stack[:len(stack)-1]   // remove the top element

			pathKey := strings.Join(current.Path, " -> ")

			// print the current checking node
			fmt.Printf("Checking: %s\n", pathKey)

			// check if the current node is the target
			if current.Title == target.Title && !pathSet[pathKey] {
				results = append(results, current.Path) // add the found path to the results
				pathSet[pathKey] = true                 // mark this path as found
				if depth < foundDepth {
					foundDepth = depth // stop the depth iteration
				}
			}

			if len(current.Path) > depth || visited[current.URL] { // if too deep or already visited, skip
				continue
			}

			visited[current.URL] = true
			urlQueue <- current.URL // mark the url as visited

			// increment the articles checked counter
			articlesChecked++

			select {
			case neighbors := <-resultQueue: // if the results are available
				for _, neighbor := range neighbors {
					if !visited[neighbor.URL] { // if the neighbor ins't visited
						neighbor.Path = append([]string(nil), current.Path...) // copy the current path
						neighbor.Path = append(neighbor.Path, neighbor.Title)  // extend the path with the neighbors
						stack = append(stack, neighbor)                        // add to the stack for further exploration

						// increment the articles traversed counter
						articlesTraversed++
					}
				}
			case err := <-errQueue: // if there's error, print it
				log.Println(err)
			}
		}
	}

	// close the channels and wait for all workers to finish
	close(urlQueue)
	wg.Wait()
	close(resultQueue)
	close(errQueue)

	elapsedTime := time.Since(startTime).Seconds()
	numberPath := len(results)
	fmt.Printf("Result: %v\n", results)
	fmt.Println("\nShortest Path(s) found:")
	for i, result := range results {
		fmt.Printf("%d.  %s\n", i+1, strings.Join(result, " -> "))
	}
	fmt.Printf("Number of shortest path(s) found: %d\n", numberPath)
	fmt.Printf("Articles Checked: %d\n", articlesChecked)
	fmt.Printf("Articles Traversed: %d\n", articlesTraversed)
	fmt.Printf("Elapsed Time: %f seconds\n", elapsedTime)

	return results, articlesChecked, articlesTraversed, numberPath, elapsedTime // return the final results and counters
}

// iterative deepening search to find the shortest path from a start node to a target node
func iterativeDeepeningShortest(start, target Node, maxDepth int) ([][]string, int, int, int, float64) {
	startTime := time.Now()
	var result [][]string                 // to store the final result
	pathSet := make(map[string]bool)      // to ensure unique paths
	urlQueue := make(chan string, 500)    // channel for the urls to be processed
	resultQueue := make(chan []Node, 500) // channel for the results
	errQueue := make(chan error, 500)     // channel for the errors

	// create multiple worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < 30; i++ { // set 30 workers
		wg.Add(1)
		go worker(urlQueue, resultQueue, errQueue, &wg) // start a new worker
	}

	// counters for articles checked and traversed
	var articlesChecked int
	var articlesTraversed int

	// flag to indicate if a path has been found
	var found bool

	// iterate through increasing depth
	for depth := 0; depth <= maxDepth && !found; depth++ {
		visited := make(map[string]bool) // track visited nodes
		var stack []Node

		// set the initial path for the start node
		start.Path = []string{start.Title}
		stack = append(stack, start)

		for len(stack) > 0 && !found {
			current := stack[len(stack)-1] // get the top of the stack
			stack = stack[:len(stack)-1]   // remove the top element

			pathKey := strings.Join(current.Path, " -> ")

			// print the current checking node
			fmt.Printf("Checking: %s\n", pathKey)

			// check if the current node is the target
			if current.Title == target.Title && !pathSet[pathKey] {
				result = append(result, current.Path) // append the found path to result
				pathSet[pathKey] = true               // mark this path as found
				found = true                          // set found flag to true
				break                                 // break the loop
			}

			if len(current.Path) > depth || visited[current.URL] { // if too deep or already visited, skip
				continue
			}

			visited[current.URL] = true
			urlQueue <- current.URL // mark the url as visited

			// increment the articles checked counter
			articlesChecked++

			select {
			case neighbors := <-resultQueue: // if the results are available
				for _, neighbor := range neighbors {
					if !visited[neighbor.URL] { // if the neighbor ins't visited
						neighbor.Path = append([]string(nil), current.Path...) // copy the current path
						neighbor.Path = append(neighbor.Path, neighbor.Title)  // extend the path with the neighbors
						stack = append(stack, neighbor)                        // add to the stack for further exploration

						// increment the articles traversed counter
						articlesTraversed++
					}
				}
			case err := <-errQueue: // if there's error, print it
				log.Println(err)
			}
		}
	}

	// close the channels and wait for all workers to finish
	close(urlQueue)
	wg.Wait()
	close(resultQueue)
	close(errQueue)

	elapsedTime := time.Since(startTime).Seconds()

	fmt.Printf("Result: %v\n", result)
	fmt.Println("\nShortest Path found:")
	if found {
		for _, path := range result {
			fmt.Printf("%s\n", strings.Join(path, " -> "))
		}
	}

	fmt.Println("Number of shortest path(s) found:", len(result))
	fmt.Printf("Articles Checked: %d\n", articlesChecked)
	fmt.Printf("Articles Traversed: %d\n", articlesTraversed)
	fmt.Printf("Elapsed Time: %f seconds\n", elapsedTime)

	return result, articlesChecked, articlesTraversed, len(result), elapsedTime // return the final result and counters
}

func replaceSpacesWithUnderscores(title string) string {
	return strings.ReplaceAll(title, " ", "_")
}

func replaceUnderscoresWithSpaces(title string) string {
	return strings.ReplaceAll(title, "_", " ")
}

// func handleSearchIDS(w http.ResponseWriter, r *http.Request) {
// 	var requestData map[string]string
// 	err := json.NewDecoder(r.Body).Decode(&requestData)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	start := requestData["start"]
// 	target := requestData["target"]

// 	fmt.Printf("Start: %s, Target: %s\n", start, target) // Debug print

// 	startNode := Node{Title: replaceUnderscoresWithSpaces(start), URL: fmt.Sprintf("https://en.wikipedia.org/wiki/%s", start), Path: []string{replaceUnderscoresWithSpaces(start)}}
// 	targetNode := Node{Title: replaceUnderscoresWithSpaces(target), URL: fmt.Sprintf("https://en.wikipedia.org/wiki/%s", target), Path: []string{replaceUnderscoresWithSpaces(target)}}

// 	results, articleChecked, articlesTraversed, numberPath, elapsedTime := iterativeDeepeningShortest(startNode, targetNode, 6)

// 	// convert results to json
// 	err = json.NewEncoder(w).Encode(Result{results, articleChecked, articlesTraversed, numberPath, elapsedTime})
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }
