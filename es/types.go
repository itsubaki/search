package es

import (
	"bytes"
	"encoding/json"
)

type CatIndex struct {
	Health       string `json:"health"`
	Status       string `json:"status"`
	Index        string `json:"index"`
	UUID         string `json:"uuid"`
	Pri          string `json:"pri"`
	Rep          string `json:"rep"`
	DocsCount    string `json:"docs.count"`
	DocsDeleted  string `json:"docs.deleted"`
	StoreSize    string `json:"store.size"`
	PriStoreSize string `json:"pri.store.size"`
	DatasetSize  string `json:"dataset.size"`
}

type SearchResult[T any] struct {
	Took     int         `json:"took"`
	TimedOut bool        `json:"timed_out"`
	Shards   Shards      `json:"_shards"`
	Hits     HitsBody[T] `json:"hits"`
}

type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

type HitsBody[T any] struct {
	Total    TotalInfo      `json:"total"`
	MaxScore float64        `json:"max_score"`
	Hits     []HitDetail[T] `json:"hits"`
}

type TotalInfo struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

type HitDetail[T any] struct {
	Index  string  `json:"_index"`
	ID     string  `json:"_id"`
	Score  float64 `json:"_score"`
	Source T       `json:"_source"`
}

type BulkIndex struct {
	Index BulkIndexMeta `json:"index"`
}

type BulkIndexMeta struct {
	Index string `json:"_index"`
	ID    int    `json:"_id"`
}

type Data[T any] struct {
	BulkIndex BulkIndex
	Source    T
}

func Bytes[T any](data []Data[T]) ([]byte, error) {
	var buf bytes.Buffer
	for _, d := range data {
		meta, err := json.Marshal(d.BulkIndex)
		if err != nil {
			return nil, err
		}

		buf.Write(meta)
		buf.WriteByte('\n')

		source, err := json.Marshal(d.Source)
		if err != nil {
			return nil, err
		}

		buf.Write(source)
		buf.WriteByte('\n')
	}

	return buf.Bytes(), nil
}

type Query struct {
	Query BoolQuery `json:"query"`
}

type BoolQuery struct {
	Bool BoolBody `json:"bool"`
}

type BoolBody struct {
	Must   []MultiMatchQuery `json:"must,omitempty"`
	Should []MultiMatchQuery `json:"should,omitempty"`
}

type MultiMatchQuery struct {
	MultiMatch MultiMatchBody `json:"multi_match"`
}

type MultiMatchBody struct {
	Query  string   `json:"query"`
	Fields []string `json:"fields"`
	Type   string   `json:"type"`
}

func (q *Query) Bytes() ([]byte, error) {
	return json.Marshal(q)
}

type Index struct {
	Settings Settings `json:"settings"`
	Mappings Mappings `json:"mappings"`
}

type Settings struct {
	Analysis Analysis `json:"analysis"`
}

type Analysis struct {
	CharFilter map[string]CharFilter `json:"char_filter"`
	Tokenizer  map[string]Tokenizer  `json:"tokenizer"`
	Filter     map[string]Filter     `json:"filter"`
	Analyzer   map[string]Analyzer   `json:"analyzer"`
}

type CharFilter struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Mode string `json:"mode"`
}

type Tokenizer struct {
	Type                 string   `json:"type"`
	Mode                 string   `json:"mode,omitempty"`
	DiscardCompoundToken bool     `json:"discard_compound_token,omitempty"`
	UserDictionaryRules  []string `json:"user_dictionary_rules,omitempty"`
	MinGram              int      `json:"min_gram,omitempty"`
	MaxGram              int      `json:"max_gram,omitempty"`
	TokenChars           []string `json:"token_chars,omitempty"`
}

type Filter struct {
	Type     string   `json:"type"`
	Lenient  bool     `json:"lenient"`
	Synonyms []string `json:"synonyms"`
}

type Analyzer struct {
	Type       string   `json:"type"`
	CharFilter []string `json:"char_filter"`
	Tokenizer  string   `json:"tokenizer"`
	Filter     []string `json:"filter"`
}

type Mappings struct {
	Properties map[string]Property `json:"properties"`
}

type Property struct {
	Type           string              `json:"type"`
	SearchAnalyzer string              `json:"search_analyzer"`
	Analyzer       string              `json:"analyzer"`
	Fields         map[string]SubField `json:"fields"`
}

type SubField struct {
	Type           string `json:"type"`
	SearchAnalyzer string `json:"search_analyzer"`
	Analyzer       string `json:"analyzer"`
}
