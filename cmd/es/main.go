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

type MyData struct {
	MyField string `json:"my_field"`
}

func main() {
	client, err := es.NewClient(
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

	if err := client.Delete([]string{
		indexName,
	}); err != nil {
		panic(err)
	}

	if err := client.Create(
		context.Background(),
		indexName,
		es.Index{
			Settings: es.Settings{
				Analysis: es.Analysis{
					CharFilter: map[string]es.CharFilter{
						"normalize": {
							Type: "icu_normalizer",
							Name: "nfkc",
							Mode: "compose",
						},
					},
					Tokenizer: map[string]es.Tokenizer{
						"ja_kuromoji_tokenizer": {
							Type:                 "kuromoji_tokenizer",
							Mode:                 "search",
							DiscardCompoundToken: true,
							UserDictionaryRules: []string{
								"東京スカイツリー,東京 スカイツリー,トウキョウ スカイツリー,カスタム名詞",
							},
						},
						"ja_ngram_tokenizer": {
							Type:    "ngram",
							MinGram: 2,
							MaxGram: 2,
							TokenChars: []string{
								"letter",
								"digit",
							},
						},
					},
					Filter: map[string]es.Filter{
						"ja_index_synonym": {
							Type:     "synonym",
							Lenient:  false,
							Synonyms: []string{},
						},
						"ja_search_synonym": {
							Type:    "synonym_graph",
							Lenient: false,
							Synonyms: []string{
								"米国, アメリカ",
								"東京大学, 東大",
							},
						},
					},
					Analyzer: map[string]es.Analyzer{
						"ja_kuromoji_index_analyzer": {
							Type:      "custom",
							Tokenizer: "ja_kuromoji_tokenizer",
							CharFilter: []string{
								"normalize",
							},
							Filter: []string{
								"kuromoji_baseform",
								"kuromoji_part_of_speech",
								"ja_index_synonym",
								"cjk_width",
								"ja_stop",
								"kuromoji_stemmer",
								"lowercase",
							},
						},
						"ja_kuromoji_search_analyzer": {
							Type:      "custom",
							Tokenizer: "ja_kuromoji_tokenizer",
							CharFilter: []string{
								"normalize",
							},
							Filter: []string{
								"kuromoji_baseform",
								"kuromoji_part_of_speech",
								"ja_search_synonym",
								"cjk_width",
								"ja_stop",
								"kuromoji_stemmer",
								"lowercase",
							},
						},
						"ja_ngram_index_analyzer": {
							Type:      "custom",
							Tokenizer: "ja_ngram_tokenizer",
							CharFilter: []string{
								"normalize",
							},
							Filter: []string{
								"lowercase",
							},
						},
						"ja_ngram_search_analyzer": {
							Type:      "custom",
							Tokenizer: "ja_ngram_tokenizer",
							CharFilter: []string{
								"normalize",
							},
							Filter: []string{
								"ja_search_synonym",
								"lowercase",
							},
						},
					},
				},
			},
			Mappings: es.Mappings{
				Properties: map[string]es.Property{
					"my_field": {
						Type:           "text",
						SearchAnalyzer: "ja_kuromoji_search_analyzer",
						Analyzer:       "ja_kuromoji_index_analyzer",
						Fields: map[string]es.Field{
							"ngram": {
								Type:           "text",
								SearchAnalyzer: "ja_ngram_search_analyzer",
								Analyzer:       "ja_ngram_index_analyzer",
							},
						},
					},
				},
			},
		},
	); err != nil {
		panic(err)
	}

	if err := es.Bulk(
		context.Background(),
		client,
		indexName,
		[]es.Data[MyData]{
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
		},
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

	for _, q := range []es.Query{
		{
			Query: es.BoolQuery{
				Bool: es.BoolBody{
					Must: []es.MultiMatchQuery{
						{
							MultiMatch: es.MultiMatchBody{
								Query: "米国",
								Type:  "phrase",
								Fields: []string{
									"my_field.ngram^1",
								},
							},
						},
					},
					Should: []es.MultiMatchQuery{
						{
							MultiMatch: es.MultiMatchBody{
								Query: "米国",
								Type:  "phrase",
								Fields: []string{
									"my_field^1",
								},
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
								Query: "東京大学",
								Type:  "phrase",
								Fields: []string{
									"my_field.ngram^1",
								},
							},
						},
					},
					Should: []es.MultiMatchQuery{
						{
							MultiMatch: es.MultiMatchBody{
								Query: "東京大学",
								Type:  "phrase",
								Fields: []string{
									"my_field^1",
								},
							},
						},
					},
				},
			},
		},
	} {
		result, err := es.Search[MyData](
			context.Background(),
			client,
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
