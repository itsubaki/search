package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/itsubaki/search/es"
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

type MyData struct {
	MyField string `json:"my_field"`
}

var data = []es.Data[MyData]{
	{
		BulkIndex: es.BulkIndex{
			Index: es.BulkIndexMeta{
				Index: indexName,
				ID:    1,
			},
		},
		Source: MyData{
			MyField: "アメリカ",
		},
	},
	{
		BulkIndex: es.BulkIndex{
			Index: es.BulkIndexMeta{
				Index: indexName,
				ID:    2,
			},
		},
		Source: MyData{
			MyField: "米国",
		},
	},
	{
		BulkIndex: es.BulkIndex{
			Index: es.BulkIndexMeta{
				Index: indexName,
				ID:    3,
			},
		},
		Source: MyData{
			MyField: "アメリカの大学",
		},
	},
	{
		BulkIndex: es.BulkIndex{
			Index: es.BulkIndexMeta{
				Index: indexName,
				ID:    4,
			},
		},
		Source: MyData{
			MyField: "東京大学",
		},
	},
	{
		BulkIndex: es.BulkIndex{
			Index: es.BulkIndexMeta{
				Index: indexName,
				ID:    5,
			},
		},
		Source: MyData{
			MyField: "帝京大学",
		},
	},
	{
		BulkIndex: es.BulkIndex{
			Index: es.BulkIndexMeta{
				Index: indexName,
				ID:    6,
			},
		},
		Source: MyData{
			MyField: "東京で夢の大学生活",
		},
	},
	{
		BulkIndex: es.BulkIndex{
			Index: es.BulkIndexMeta{
				Index: indexName,
				ID:    7,
			},
		},
		Source: MyData{
			MyField: "東京大学で夢の生活",
		},
	},
	{
		BulkIndex: es.BulkIndex{
			Index: es.BulkIndexMeta{
				Index: indexName,
				ID:    8,
			},
		},
		Source: MyData{
			MyField: "東大で夢の生活",
		},
	},
	{
		BulkIndex: es.BulkIndex{
			Index: es.BulkIndexMeta{
				Index: indexName,
				ID:    9,
			},
		},
		Source: MyData{
			MyField: "首都圏の大学 東京",
		},
	},
}

var query = []es.Query{
	{
		Query: es.BoolQuery{
			Bool: es.BoolBody{
				Must: []es.MultiMatchQuery{
					{
						MultiMatch: es.MultiMatchBody{
							Query:  "米国",
							Fields: []string{"my_field.ngram^1"},
							Type:   "phrase",
						},
					},
				},
				Should: []es.MultiMatchQuery{
					{
						MultiMatch: es.MultiMatchBody{
							Query:  "米国",
							Fields: []string{"my_field^1"},
							Type:   "phrase",
						},
					},
				},
			},
		},
	},
	{
		Query: es.BoolQuery{
			Bool: es.BoolBody{
				Must: []es.MultiMatchQuery{
					{
						MultiMatch: es.MultiMatchBody{
							Query:  "東京大学",
							Fields: []string{"my_field.ngram^1"},
							Type:   "phrase",
						},
					},
				},
				Should: []es.MultiMatchQuery{
					{
						MultiMatch: es.MultiMatchBody{
							Query:  "東京大学",
							Fields: []string{"my_field^1"},
							Type:   "phrase",
						},
					},
				},
			},
		},
	},
}

func main() {
	client, err := es.NewClient[MyData](
		[]string{address},
		username,
		password,
	)
	if err != nil {
		panic(err)
	}

	if err := client.Ping(); err != nil {
		panic(err)
	}

	settings, err := json.Marshal(index)
	if err != nil {
		panic(err)
	}

	if err := client.Delete([]string{
		indexName,
	}); err != nil {
		panic(err)
	}

	if err := client.Create(
		context.Background(),
		indexName,
		settings,
	); err != nil {
		panic(err)
	}

	if err := client.Bulk(
		context.Background(),
		indexName,
		data,
	); err != nil {
		panic(err)
	}

	for {
		cnt, err := client.Count(indexName)
		if err != nil {
			panic(err)
		}

		if cnt > 0 {
			break
		}

		time.Sleep(1 * time.Second)
	}

	list, err := client.CatIndex()
	if err != nil {
		panic(err)
	}

	for _, v := range list {
		bytes, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(bytes))
	}

	for _, q := range query {
		result, err := client.Search(
			context.Background(),
			indexName,
			q,
		)
		if err != nil {
			panic(err)
		}

		bytes, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(bytes))
	}
}
