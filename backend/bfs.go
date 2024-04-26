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
        if encountered[value] == false {
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
	
	startTime := time.Now()
	bfsMulti(startURL, targetTitle, hrefs, &mu, &articlesChecked, &results)
	elapsedTime := time.Since(startTime).Seconds()
	
	for i, element := range results {
		results[i] = removeDuplicates(element)
	}	
	results = removeDuplicateLists(results)
	numberPath := len(results)

	return results, articlesChecked, articlesTraversed, numberPath, elapsedTime
} 

func bfsMulti(startURL, targetTitle string, hrefs []string, mu *sync.Mutex, passed *int, results *[][]string) {
	depth := 4
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

				// fmt.Println(strings.TrimSpace(firstPathSplit[0]))
				
				content := doc.Find("#mw-content-text")
				content.Find("a").Each(func(i int, s *goquery.Selection) {
					href, exists := s.Attr("href")
					if exists && strings.HasPrefix(href, "/wiki/") {
						title := strings.TrimSpace(firstPathSplit[0])
							if title == targetTitle {
								result := append(page.Path, title)
								*results = append(*results, result)
								// result := strings.Join(page.Path, " -> ")
								// fmt.Printf("Path: %s -> %s\n", result, title)
								return
							} else if href != "/wiki/Main_Page" && !isInList(href, hrefs) {
								hrefs = append(hrefs, href)
								mu.Lock()
								*passed++
								queue = append(queue, Page{URL: "https://en.wikipedia.org" + href, Path: page.Path, Depth: page.Depth - 1})
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
	
	startTime := time.Now()
	bfsSingle(startURL, targetTitle, hrefs, &mu, &articlesChecked, &found, &results)
	elapsedTime := time.Since(startTime).Seconds()
	
	for i, element := range results {
		results[i] = removeDuplicates(element)
	}	
	results = removeDuplicateLists(results)
	numberPath := len(results)

	return results, articlesChecked, articlesTraversed, numberPath, elapsedTime
} 

func bfsSingle(startURL, targetTitle string, hrefs []string, mu *sync.Mutex, passed *int, found *bool, results *[][]string) {
	depth := 2
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
				firstPathSplit := strings.Split(firstPath, " - ")
				if len(firstPathSplit) > 0 {
					page.Path = append(page.Path, strings.TrimSpace(firstPathSplit[0]))
				}

				content := doc.Find("#mw-content-text")
				content.Find("a").Each(func(i int, s *goquery.Selection) {
					href, exists := s.Attr("href")
					if exists && strings.HasPrefix(href, "/wiki/") {
						title := firstPathSplit[0]
						// fmt.Printf("checking %s\n", title)
						if title == targetTitle {
							result := append(page.Path, title)
							*results = append(*results, result)
							*found = true
							// fmt.Printf("Path: %s -> %s\n", result, title)
							// fmt.Printf("Path Length: %d\n", page.Depth)
							// fmt.Printf("Checked Article: %d\n", *passed)
							return
						} else if href != "/wiki/Main_Page" && !isInList(href, hrefs) {
							hrefs = append(hrefs, href)
							mu.Lock()
							*passed++
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