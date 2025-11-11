package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

const (
	user = "admin"
	pswd = "xuYz3_cAXYh7"
	addr = "https://localhost:9200"
)

func main() {
	client, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Addresses: []string{
			addr,
		},
		Username: user,
		Password: pswd,
	})
	if err != nil {
		panic(err)
	}

	resp, err := client.Info()
	if err != nil {
		panic(err)
	}
	fmt.Println(MustRead(resp.Body))

	{
		resp, err := opensearchapi.IndicesDeleteRequest{
			Index: []string{"go-test-index1"},
		}.Do(context.Background(), client)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println(string(MustRead(resp.Body)))
	}

	{
		resp, err := opensearchapi.IndicesCreateRequest{
			Index: "go-test-index1",
			Body:  strings.NewReader(`{"settings": {"index": {"number_of_shards": 1, "number_of_replicas": 0}}}`),
		}.Do(context.Background(), client)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println(string(MustRead(resp.Body)))
	}

	{
		resp, err := opensearchapi.IndexRequest{
			Index:      "go-test-index1",
			DocumentID: "1",
			Body:       strings.NewReader(`{"title": "Moneyball", "director": "Bennett Miller", "year": "2011"}`),
		}.Do(context.Background(), client)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println(MustRead(resp.Body))
	}

	{
		data := `
	{ "index" : { "_index" : "go-test-index1", "_id" : "2" } }
	{ "title" : "Interstellar", "director" : "Christopher Nolan", "year" : "2014"}
	{ "create" : { "_index" : "go-test-index1", "_id" : "3" } }
	{ "title" : "Star Trek Beyond", "director" : "Justin Lin", "year" : "2015"}
	{ "update" : {"_id" : "3", "_index" : "go-test-index1" } }
	{ "doc" : {"year" : "2016"} }
	` + "\n"

		resp, err := client.Bulk(strings.NewReader(data))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println(MustRead(resp.Body))
	}

	{
		content := strings.NewReader(`{"size": 5, "query": { "multi_match": { "query": "miller", "fields": ["title^2", "director"]}}}`)
		resp, err := opensearchapi.SearchRequest{
			Index: []string{"go-test-index1"},
			Body:  content,
		}.Do(context.Background(), client)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println(MustRead(resp.Body))
	}
}

func MustRead(r io.Reader) string {
	bytes, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}
