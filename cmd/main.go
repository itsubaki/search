package main

import (
	"context"
	"fmt"
	"os"

	"github.com/itsubaki/search/osr"
)

const (
	username  = "admin"
	password  = "xuYz3_cAXYh7"
	addr      = "https://localhost:9200"
	indexName = "my_fulltext_search"
)

var (
	index = func() []byte {
		index, err := os.ReadFile("testdata/index.json")
		if err != nil {
			panic(err)
		}

		return index
	}()

	bulk = func() []osr.Data[MyData] {
		data, err := os.ReadFile("testdata/data.jsonl")
		if err != nil {
			panic(err)
		}

		bulk, err := osr.Read[MyData](data)
		if err != nil {
			panic(err)
		}

		return bulk
	}()

	query = func() [][]byte {
		us, err := os.ReadFile("testdata/query_us.json")
		if err != nil {
			panic(err)
		}

		jp, err := os.ReadFile("testdata/query_jp.json")
		if err != nil {
			panic(err)
		}

		return [][]byte{
			us,
			jp,
		}
	}()
)

type MyData struct {
	MyField       string    `json:"my_field"`
	MyDenseVector []float32 `json:"my_dense_vector,omitempty"`
}

func main() {
	client, err := osr.NewClient(
		[]string{addr},
		username,
		password,
	)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	if err := client.Delete(ctx, []string{indexName}); err != nil {
		fmt.Println(err)
	}

	if err := client.Create(ctx, indexName, index); err != nil {
		panic(err)
	}

	if err := osr.Bulk(ctx, client, bulk); err != nil {
		panic(err)
	}

	if err := client.Refresh(ctx, []string{indexName}); err != nil {
		panic(err)
	}

	for _, q := range query {
		results, err := osr.Search[MyData](
			ctx,
			client,
			[]string{indexName},
			q,
		)
		if err != nil {
			panic(err)
		}

		for _, hit := range results.Hits.Hits {
			fmt.Printf("%2.4f, %+v\n", hit.Score, hit.Source)
		}
		fmt.Println()
	}
}
