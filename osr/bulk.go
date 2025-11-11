package osr

import (
	"bytes"
	"encoding/json"
)

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
