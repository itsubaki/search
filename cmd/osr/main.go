package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/itsubaki/search/osr"
)

const (
	username = "admin"
	password = "xuYz3_cAXYh7"
	addr     = "https://localhost:9200"
	index    = "go-test-index1"
)

func main() {
	client, err := osr.NewClient(
		[]string{addr},
		username,
		password,
	)
	if err != nil {
		panic(err)
	}

	resp, err := client.CatIndex()
	if err != nil {
		panic(err)
	}

	for _, r := range resp {
		fmt.Printf("%+v\n", r)
	}
	fmt.Println()

	ctx := context.Background()
	if err := client.Delete(ctx, []string{index}); err != nil {
		fmt.Println(err)
	}

	if err := client.Create(ctx, index, strings.NewReader(
		`{"settings": {"index": {"number_of_shards": 1, "number_of_replicas": 0}}}`,
	)); err != nil {
		panic(err)
	}

	type Movie struct {
		Title    string `json:"title"`
		Director string `json:"director"`
		Year     string `json:"year"`
	}

	if err := osr.Bulk(ctx, client, []osr.Data[Movie]{
		{
			BulkIndex: osr.BulkIndex{
				Index: osr.BulkIndexMeta{
					Index: index,
					ID:    1,
				},
			},
			Source: Movie{
				Title:    "Moneyball",
				Director: "Bennett Miller",
				Year:     "2011",
			},
		},
		{
			BulkIndex: osr.BulkIndex{
				Index: osr.BulkIndexMeta{
					Index: index,
					ID:    2,
				},
			},
			Source: Movie{
				Title:    "Interstellar",
				Director: "Christopher Nolan",
				Year:     "2014",
			},
		},
		{
			BulkIndex: osr.BulkIndex{
				Index: osr.BulkIndexMeta{
					Index: index,
					ID:    3,
				},
			},
			Source: Movie{
				Title:    "Star Trek Beyond",
				Director: "Justin Lin",
				Year:     "2016",
			},
		},
	}); err != nil {
		panic(err)
	}

	if err := client.Refresh(ctx, []string{index}); err != nil {
		panic(err)
	}

	results, err := osr.Search[Movie](
		ctx,
		client,
		[]string{index},
		strings.NewReader(
			`{"size": 5, "query": { "multi_match": { "query": "miller", "fields": ["title^2", "director"]}}}`,
		),
	)
	if err != nil {
		panic(err)
	}

	for _, hit := range results.Hits.Hits {
		fmt.Printf("%+v\n", hit.Source)
	}
}
