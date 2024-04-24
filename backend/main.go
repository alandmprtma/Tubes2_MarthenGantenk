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

/* Fungsi Menampilkan Hasil Pencarian Dari Wikipedia API */
func handleWikipediaRequest(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Query parameter is required", http.StatusBadRequest)
		return
	}

	wikipediaURL := "https://en.wikipedia.org/w/api.php?action=query&format=json&list=search&srsearch=" + url.QueryEscape(query)
	response, err := http.Get(wikipediaURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	var data interface{}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
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
