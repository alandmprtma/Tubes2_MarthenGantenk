package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
        r := mux.NewRouter()
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