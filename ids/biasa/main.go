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

func isInList(node Node, list []Node) bool {
	for _, v := range list {
		if v.URL == node.URL {
			return true
		}
	}
	return false
}

func FindLink(startNode, targetNode Node, depth int, visited []Node) {
	if depth < 0 {
		return
	}

	res, err := http.Get(startNode.URL)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
		return
	}

	bodyContent := doc.Find("#bodyContent")
	if bodyContent.Length() == 0 {
		log.Println("bodyContent not found")
		return
	}

	bodyContent.Find("p").Each(func(i int, p *goquery.Selection) {
		p.Find("a").Each(func(j int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if exists && strings.HasPrefix(href, "/wiki/") {
				title := s.Text()
				fmt.Printf("Checking: %s\n", title)
				if title == targetNode.Title {
					fmt.Printf("Found: %s\n", title)
					path := append(startNode.Path, title)
					fmt.Printf("Path: %s\n", strings.Join(path, " -> "))
					panic(path)
				}
				nextNode := Node{Title: title, URL: "https://en.wikipedia.org" + href, Path: append(startNode.Path, title)}
				if href != "/wiki/Main_Page" && !isInList(nextNode, visited) {
					visited = append(visited, nextNode)
					FindLink(nextNode, targetNode, depth-1, visited)
				}
			}
		})
	})
}

func main() {
	startNode := Node{Title: "Mike Tyson", URL: "https://en.wikipedia.org/wiki/Mike_Tyson", Path: []string{"Mike Tyson"}}
	targetNode := Node{Title: "boxing", URL: "", Path: []string{}}
	depth := 3
	var visited []Node

	startTime := time.Now()
	defer func() {
		if r := recover(); r != nil {
			path := r.([]string)
			fmt.Printf("Shortest path: %s\n", strings.Join(path, " -> "))
			elapsedTime := time.Since(startTime)
			fmt.Printf("Elapsed time: %v\n", elapsedTime)
		}
	}()
	FindLink(startNode, targetNode, depth, visited)
}
