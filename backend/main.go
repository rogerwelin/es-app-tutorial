package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/olivere/elastic"
)

const (
	elasticIndexName = "cartoon"
	elasticTypeName  = "cartoon"
	elasticURL       = "http://localhost:9200"
)

type AnimeDocument struct {
	Name     string   `json:"name"`
	Genre    []string `json:"genre"`
	Type     string   `json:"type"`
	Episodes string   `json:"episodes"`
	Rating   string   `json:"rating"`
}

func main() {
	client, err := elastic.NewClient()
	if err != nil {
		panic(err)
	}

	// search
	termQuery := elastic.NewTermQuery("name", "berserk")
	searchResult, err := client.Search().
		Index("cartoon").
		Query(termQuery).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	// iterate over search result
	if searchResult.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d results\n", searchResult.Hits.TotalHits)
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a AnimeDocument (could also be just a map[string]interface{}).
			var a AnimeDocument
			err := json.Unmarshal(*hit.Source, &a)
			if err != nil {
				fmt.Println("Deserialization failed")
			}

			fmt.Printf("Result --> %s: %s %s %s %s\n", a.Name, a.Genre, a.Type, a.Episodes, a.Rating)
		}
	} else {
		// No hits
		fmt.Print("Found no results\n")
	}

}
