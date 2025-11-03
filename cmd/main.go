package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
)

var (
	indexName = "my_fulltext_search"
	username  = "elastic"
	password  = os.Getenv("ES_PASSWORD")
	address   = "https://localhost:9200"
)

var index = map[string]any{
	"settings": map[string]any{
		"analysis": map[string]any{
			"char_filter": map[string]any{
				"normalize": map[string]any{
					"type": "icu_normalizer",
					"name": "nfkc",
					"mode": "compose",
				},
			},
			"tokenizer": map[string]any{
				"ja_kuromoji_tokenizer": map[string]any{
					"type":                   "kuromoji_tokenizer",
					"mode":                   "search",
					"discard_compound_token": true,
					"user_dictionary_rules": []string{
						"東京スカイツリー,東京 スカイツリー,トウキョウ スカイツリー,カスタム名詞",
					},
				},
				"ja_ngram_tokenizer": map[string]any{
					"type":     "ngram",
					"min_gram": 2,
					"max_gram": 2,
					"token_chars": []string{
						"letter",
						"digit",
					},
				},
			},
			"filter": map[string]any{
				"ja_index_synonym": map[string]any{
					"type":     "synonym",
					"lenient":  false,
					"synonyms": []string{},
				},
				"ja_search_synonym": map[string]any{
					"type":    "synonym_graph",
					"lenient": false,
					"synonyms": []string{
						"米国, アメリカ",
						"東京大学, 東大",
					},
				},
			},
			"analyzer": map[string]any{
				"ja_kuromoji_index_analyzer": map[string]any{
					"type":        "custom",
					"char_filter": []string{"normalize"},
					"tokenizer":   "ja_kuromoji_tokenizer",
					"filter": []string{
						"kuromoji_baseform",
						"kuromoji_part_of_speech",
						"ja_index_synonym",
						"cjk_width",
						"ja_stop",
						"kuromoji_stemmer",
						"lowercase",
					},
				},
				"ja_kuromoji_search_analyzer": map[string]any{
					"type":        "custom",
					"char_filter": []string{"normalize"},
					"tokenizer":   "ja_kuromoji_tokenizer",
					"filter": []string{
						"kuromoji_baseform",
						"kuromoji_part_of_speech",
						"ja_search_synonym",
						"cjk_width",
						"ja_stop",
						"kuromoji_stemmer",
						"lowercase",
					},
				},
				"ja_ngram_index_analyzer": map[string]any{
					"type":        "custom",
					"char_filter": []string{"normalize"},
					"tokenizer":   "ja_ngram_tokenizer",
					"filter":      []string{"lowercase"},
				},
				"ja_ngram_search_analyzer": map[string]any{
					"type":        "custom",
					"char_filter": []string{"normalize"},
					"tokenizer":   "ja_ngram_tokenizer",
					"filter": []string{
						"ja_search_synonym",
						"lowercase",
					},
				},
			},
		},
	},
	"mappings": map[string]any{
		"properties": map[string]any{
			"my_field": map[string]any{
				"type":            "text",
				"search_analyzer": "ja_kuromoji_search_analyzer",
				"analyzer":        "ja_kuromoji_index_analyzer",
				"fields": map[string]any{
					"ngram": map[string]any{
						"type":            "text",
						"search_analyzer": "ja_ngram_search_analyzer",
						"analyzer":        "ja_ngram_index_analyzer",
					},
				},
			},
		},
	},
}

var data = `
{"index": {"_index": "my_fulltext_search", "_id": 1}}
{"my_field": "アメリカ"}
{"index": {"_index": "my_fulltext_search", "_id": 2}}
{"my_field": "米国"}
{"index": {"_index": "my_fulltext_search", "_id": 3}}
{"my_field": "アメリカの大学"}
{"index": {"_index": "my_fulltext_search", "_id": 4}}
{"my_field": "東京大学"}
{"index": {"_index": "my_fulltext_search", "_id": 5}}
{"my_field": "帝京大学"}
{"index": {"_index": "my_fulltext_search", "_id": 6}}
{"my_field": "東京で夢の大学生活"}
{"index": {"_index": "my_fulltext_search", "_id": 7}}
{"my_field": "東京大学で夢の生活"}
{"index": {"_index": "my_fulltext_search", "_id": 8}}
{"my_field": "東大で夢の生活"}
{"index": {"_index": "my_fulltext_search", "_id": 9}}
{"my_field": "首都圏の大学 東京"}
`

var query = []string{
	`
{
  "query": {
    "bool": {
      "must": [
        {
          "multi_match": {
            "query": "米国",
            "fields": ["my_field.ngram^1"],
            "type": "phrase"
          }
        }
      ],
      "should": [
        {
          "multi_match": {
            "query": "米国",
            "fields": ["my_field^1"],
            "type": "phrase"
          }
        }
      ]
    }
  }
}
`,
	`
{
  "query": {
    "bool": {
      "must": [
        {
          "multi_match": {
            "query": "東京大学",
            "fields": [
              "my_field.ngram^1"
            ],
            "type": "phrase"
          }
        }
      ],
      "should": [
        {
          "multi_match": {
            "query": "東京大学",
            "fields": [
              "my_field^1"
            ],
            "type": "phrase"
          }
        }
      ]
    }
  }
}
`,
}

func main() {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			address,
		},
		Username: username,
		Password: password,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	pong, err := es.Ping()
	if err != nil {
		panic(err)
	}
	defer pong.Body.Close()
	fmt.Println(pong)

	{
		data, err := json.Marshal(index)
		if err != nil {
			panic(err)
		}

		if _, err = es.Indices.Delete([]string{indexName}); err != nil {
			panic(err)
		}

		res, err := es.Indices.Create(
			indexName,
			es.Indices.Create.WithBody(bytes.NewReader(data)),
			es.Indices.Create.WithContext(context.Background()),
		)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		if res.IsError() {
			panic(res.String())
		}

		fmt.Println("✅ Successfully created index")
	}

	{
		res, err := es.Bulk(
			bytes.NewReader([]byte(data)),
			es.Bulk.WithContext(context.Background()),
			es.Bulk.WithIndex(indexName),
		)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		if res.IsError() {
			panic(res.String())
		}

		fmt.Println("✅ Bulk indexing completed")

		for {
			res, err := es.Count(
				es.Count.WithIndex(indexName),
			)
			if err != nil {
				panic(err)
			}

			type Response struct {
				Count  int `json:"count"`
				Shards struct {
					Total      int `json:"total"`
					Successful int `json:"successful"`
					Skipped    int `json:"skipped"`
					Failed     int `json:"failed"`
				} `json:"_shards"`
			}

			var count Response
			if err := json.NewDecoder(res.Body).Decode(&count); err != nil {
				panic(err)
			}

			if count.Count > 0 {
				fmt.Println("✅ Document exists, breaking loop")
				break
			}

			fmt.Println("⏳ Waiting for documents to be indexed...")
			time.Sleep(1 * time.Second)
		}

	}

	{
		res, err := es.Cat.Indices(
			es.Cat.Indices.WithFormat("json"),
			es.Cat.Indices.WithPretty(),
		)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(body))
	}

	{
		for _, q := range query {
			res, err := es.Search(
				es.Search.WithContext(context.Background()),
				es.Search.WithIndex(indexName),
				es.Search.WithBody(bytes.NewReader([]byte(q))),
				es.Search.WithTrackTotalHits(true),
			)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()

			if res.IsError() {
				panic(res.String())
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}

			fmt.Println(string(body))
		}
	}
}
