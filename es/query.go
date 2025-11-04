package es

import "encoding/json"

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
