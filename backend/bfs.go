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

func removeDuplicateLists(lists [][]string) [][]string {
	encountered := map[string]bool{}
	result := [][]string{}

	for _, list := range lists {
		listString := fmt.Sprintf("%v", list)
		if !encountered[listString] {
			encountered[listString] = true
			result = append(result, list)
		}
	}

	return result
}

func removeDuplicates(list []string) []string {
    encountered := map[string]bool{}
    result := []string{}

    for _, value := range list {
        if !encountered[value] {
            encountered[value] = true
            result = append(result, value)
        }
    }

    return result
}

func isInList(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func bfsMultiCall(startURL string, targetTitle string) ([][]string, int, int, int, float64) {
	var results [][]string
	var articlesChecked int
	articlesTraversed := 0 // gaada disuruh di spek (?)
	
	var hrefs []string
	var mu sync.Mutex
	found := false
	depth := 2
	multiDepth := 10
	
	startTime := time.Now()
	bfsMulti(startURL, targetTitle, depth, hrefs, &mu, &articlesChecked, &articlesTraversed, &found, &results, &multiDepth)
	elapsedTime := time.Since(startTime).Seconds()
	
	for i, element := range results {
		results[i] = removeDuplicates(element)
	}	
	results = removeDuplicateLists(results)
	numberPath := len(results)

	return results, articlesChecked, articlesTraversed, numberPath, elapsedTime
} 

func bfsMulti(startURL, targetTitle string, depth int, hrefs []string, mu *sync.Mutex, checked *int, traversed *int, found *bool, results *[][]string, multiDepth *int) {
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

				if page.Depth > *multiDepth {
					return
				}

				res, err := http.Get(page.URL)
				if err != nil {
					log.Print(err)
				}
				defer res.Body.Close()
				if res.StatusCode != 200 {
					fmt.Println(page.URL)
					log.Printf("status code error: %d %s", res.StatusCode, res.Status)
					return
				}

				doc, err := goquery.NewDocumentFromReader(res.Body)
				if err != nil {
					log.Print(err)
				}

				firstPath := doc.Find("title").Text()
				firstPathSplit := strings.Split(firstPath, "-")
				if len(firstPathSplit) > 0 {
					*checked += page.Depth
					page.Path = append(page.Path, strings.TrimSpace(firstPathSplit[0]))
				}

				content := doc.Find("#mw-content-text")
				content.Find("a").Each(func(i int, s *goquery.Selection) {
					href, exists := s.Attr("href")
					if exists && strings.HasPrefix(href, "/wiki/") && !strings.Contains(href, ":") {
						// Extract the title from the URL
						title := strings.TrimPrefix(href, "/wiki/")
						title = strings.ReplaceAll(title, "_", " ") // Replace underscores with spaces
						fmt.Printf("checking %v\n", page.Path)
						if title == targetTitle {
							result := append(page.Path, title)
							*results = append(*results, result)
							*found = true
							// fmt.Printf("Path: %s -> %s\n", result, title)
							// fmt.Printf("Path Length: %d\n", page.Depth)
							// fmt.Printf("Checked Article: %d\n", *checked)
							*multiDepth = page.Depth
							return
						} else if href != "/wiki/Main_Page" && !isInList(href, hrefs) {
							hrefs = append(hrefs, href)
							mu.Lock()
							*traversed++
							queue = append(queue, Page{URL: "https://en.wikipedia.org" + href, Path: page.Path, Depth: page.Depth + 1})
							mu.Unlock()
						}
					}
				})
			}(currentPage)
		}
		wg.Wait()
	}
}

func bfsSingleCall(startURL string, targetTitle string) ([][]string, int, int, int, float64) {
	var results [][]string
	var articlesChecked int
	articlesTraversed := 0 // gaada disuruh di spek (?)
	
	var hrefs []string
	var mu sync.Mutex
	found := false
	depth := 2
	
	startTime := time.Now()
	bfsSingle(startURL, targetTitle, depth, hrefs, &mu, &articlesChecked, &articlesTraversed, &found, &results)
	elapsedTime := time.Since(startTime).Seconds()
	
	results[0] = removeDuplicates(results[0])

	numberPath := 1

	return results[0:1], articlesChecked, articlesTraversed, numberPath, elapsedTime
} 

func bfsSingle(startURL, targetTitle string, depth int, hrefs []string, mu *sync.Mutex, checked *int, traversed *int, found *bool, results *[][]string) {
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
				
				if *found {
					return
				}

				res, err := http.Get(page.URL)
				if err != nil {
					log.Print(err)
				}
				defer res.Body.Close()
				if res.StatusCode != 200 {
					log.Printf("status code error: %d %s", res.StatusCode, res.Status)
					return
				}

				doc, err := goquery.NewDocumentFromReader(res.Body)
				if err != nil {
					log.Print(err)
				}
				
				firstPath := doc.Find("title").Text()
				firstPathSplit := strings.Split(firstPath, "-")
				if len(firstPathSplit) > 0 {
					*checked += page.Depth
					page.Path = append(page.Path, strings.TrimSpace(firstPathSplit[0]))
				}
				
				content := doc.Find("#mw-content-text")
				content.Find("a").Each(func(i int, s *goquery.Selection) {
					href, exists := s.Attr("href")
					if exists && strings.HasPrefix(href, "/wiki/") && !strings.Contains(href, ":") {
						// Extract the title from the URL
						title := strings.TrimPrefix(href, "/wiki/")
						title = strings.ReplaceAll(title, "_", " ") // Replace underscores with spaces
						fmt.Printf("checking %v\n", page.Path)
						if title == targetTitle {
							result := append(page.Path, title)
							*results = append(*results, result)
							*found = true
							// fmt.Printf("Path: %s -> %s\n", result, title)
							// fmt.Printf("Path Length: %d\n", page.Depth)
							// fmt.Printf("Checked Article: %d\n", *checked)
							return
							} else if href != "/wiki/Main_Page" && !isInList(href, hrefs) {
								hrefs = append(hrefs, href)
								mu.Lock()
								*traversed++
								queue = append(queue, Page{URL: "https://en.wikipedia.org" + href, Path: page.Path, Depth: page.Depth + 1})
								mu.Unlock()
						}
					}
				})
			}(currentPage)
		}
		wg.Wait()
	}
}

// func main() {
// 	startURL := "https://en.wikipedia.org/wiki/Bandung"
// 	targetTitle := "Japan"
// 	depth := 3
// 	var hrefs []string
// 	var mu sync.Mutex
// 	start := time.Now()
// 	FindLinkMulti(startURL, targetTitle, depth, hrefs, &mu)
// 	elapsed := time.Since(start)
// 	fmt.Printf("Waktu yang dibutuhkan: %s\n", elapsed)
// }