package main

import (
	"context"
	"fmt"

	"github.com/olivere/elastic"
)

const (
	elasticIndexName = "cartoon"
	elasticTypeName  = "cartoon"
	elasticURL       = "http://localhost:9200"
)

type Document struct {
	Name string `json:"name"`
}

func main() {
	client, err := elastic.NewClient()
	if err != nil {
		panic(err)
	}

	// search
	termQuery := elastic.NewTermQuery("name", "Berserk")
	searchResult, err := client.Search().
		Index("cartoon").
		Query(termQuery).
		From(0).Size(10).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

}
