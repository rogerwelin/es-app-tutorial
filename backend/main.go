package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/olivere/elastic"
)

var (
	wait   time.Duration
	client *elastic.Client
)

const (
	elasticIndex = "cartoon"
	elasticURL   = "http://localhost:9200"
)

type AnimeDocument struct {
	Name     string   `json:"name"`
	Genre    []string `json:"genre"`
	Type     string   `json:"type"`
	Episodes string   `json:"episodes"`
	Rating   string   `json:"rating"`
}

type AnimeSearchResponse struct {
	Time      int64           `json:"time"`
	Hits      int64           `json:"hits"`
	Documents []AnimeDocument `json:"documents"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "home\n")
}

func searchElastic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	searchQuery := vars["searchTerm"]

	// search
	termQuery := elastic.NewTermQuery("name", searchQuery)

	searchResult, err := client.Search().
		Index(elasticIndex).
		Query(termQuery).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	log.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	// iterate over search result
	var res AnimeSearchResponse
	var docs []AnimeDocument
	if searchResult.Hits.TotalHits > 0 {
		log.Printf("Found a total of %d results\n", searchResult.Hits.TotalHits)
		for _, hit := range searchResult.Hits.Hits {
			// Deserialize hit.Source into a AnimeDocument (could also be just a map[string]interface{}).
			var a AnimeDocument
			err := json.Unmarshal(*hit.Source, &a)
			if err != nil {
				http.Error(w, "Could not deserialize json", http.StatusInternalServerError)
			}
			docs = append(docs, a)
		}
		res = AnimeSearchResponse{
			Time:      searchResult.TookInMillis,
			Hits:      searchResult.Hits.TotalHits,
			Documents: docs,
		}
	} else {
		res = AnimeSearchResponse{
			Time: searchResult.TookInMillis,
			Hits: searchResult.Hits.TotalHits,
		}
	}
	jsonResp, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Could not serialize json", http.StatusInternalServerError)

	}
	fmt.Fprintln(w, string(jsonResp))
}

func main() {
	var err error
	client, err = elastic.NewClient(elastic.SetURL(elasticURL))
	if err != nil {
		panic(err)
	}

	wait = time.Second * 5

	router := mux.NewRouter()
	router.HandleFunc("/", rootHandler).Methods("GET")
	router.HandleFunc("/search/{searchTerm}", searchElastic).Methods("GET")

	srv := &http.Server{
		Addr:         "0.0.0.0:4000",
		WriteTimeout: wait,
		ReadTimeout:  wait,
		Handler:      router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// graceful shutdown at SIGINT
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down server...")
	os.Exit(0)
}
