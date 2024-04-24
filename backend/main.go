package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type SearchResult struct {
    Title     string `json:"title"`
	PageID    int64  `json:"pageid"` 
    Thumbnail struct {
        Source string `json:"source"`
    } `json:"thumbnail"`
}

type SearchResponse struct {
	Query struct {
		Search []SearchResult `json:"search"`
	} `json:"query"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/search", handleSearchRequest).Methods("POST")
	// r.HandleFunc("/scrape", handleScrape).Methods("GET")
	r.HandleFunc("/api/wikipedia", handleWikipediaRequest).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	fmt.Println("server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

// Fungsi untuk fetching gambar
// Function to fetch detailed information for a given page ID
func fetchPageDetails(pageID int64) (string, error) {
    detailsURL := fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=query&format=json&prop=pageimages&pageids=%d&pithumbsize=500", pageID)
    resp, err := http.Get(detailsURL)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var result struct {
        Query struct {
            Pages map[string]struct {
                Thumbnail struct {
                    Source string `json:"source"`
                } `json:"thumbnail"`
            } `json:"pages"`
        } `json:"query"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", err
    }

    for _, page := range result.Query.Pages {
        if page.Thumbnail.Source != "" {
            return page.Thumbnail.Source, nil
        }
    }

    return "", fmt.Errorf("no image found for page ID: %d", pageID)
}

/* Fungsi Menampilkan Hasil Pencarian Dari Wikipedia API */
/* Fungsi Menampilkan Hasil Pencarian Dari Wikipedia API dengan hanya judul dan link gambar */
func handleWikipediaRequest(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Query parameter is required", http.StatusBadRequest)
		return
	}

	// Update the URL to fetch necessary data only
	wikipediaURL := "https://en.wikipedia.org/w/api.php?action=query&format=json&list=search&srsearch=" + url.QueryEscape(query) + "&srlimit=10&srprop=snippet|thumbnail"

	response, err := http.Get(wikipediaURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	var searchResponse SearchResponse
	err = json.NewDecoder(response.Body).Decode(&searchResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Adjust the data structure to include only the title and thumbnail
	var searchResults []map[string]interface{}
	for _, result := range searchResponse.Query.Search {
		searchResult := map[string]interface{}{
			"title": result.Title,
		}
		imageLink, err := fetchPageDetails(result.PageID) // Now correctly using PageID
		if err != nil {
			log.Printf("Failed to fetch image for page ID %d: %s", result.PageID, err)
			searchResult["thumbnail"] = ""
		} else {
			searchResult["thumbnail"] = imageLink
		}
		searchResults = append(searchResults, searchResult)
	}
	

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(searchResults)
}

func handleSearchRequest(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Paths             [][]string `json:"paths"`
		NumberOfPaths     int        `json:"numberOfPaths"`
		ArticlesChecked   int        `json:"articlesChecked"`
		ArticlesTraversed int        `json:"articlesTraversed"`
		ElapsedTime       float64    `json:"elapsedTime"`
	}

	var body map[string]string
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	start := body["start"]
	target := body["target"]
	// convert start and target to Node
	startNode := Node{Title: start, URL: "https://en.wikipedia.org/wiki/" + replaceSpacesWithUnderscores(start)}
	targetNode := Node{Title: target, URL: "https://en.wikipedia.org/wiki/" + replaceSpacesWithUnderscores(target)}
	paths, articlesChecked, articlesTraversed, numberPath, elapsedTime := iterativeDeepeningAll(startNode, targetNode, 5) // you can adjust the maxDepth as needed

	response := Response{
		Paths:             paths,
		NumberOfPaths:     numberPath,
		ArticlesChecked:   articlesChecked,
		ArticlesTraversed: articlesTraversed,
		ElapsedTime:       elapsedTime,
	}

	json.NewEncoder(w).Encode(response)
}
