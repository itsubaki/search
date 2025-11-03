package es

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/elastic/go-elasticsearch/v9"
)

type Client[T any] struct {
	client *elasticsearch.Client
}

func NewClient[T any](address []string, username, password string) (*Client[T], error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: address,
		Username:  username,
		Password:  password,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &Client[T]{
		client: es,
	}, nil
}

func (c *Client[T]) Ping() error {
	if _, err := c.client.Ping(); err != nil {
		return err
	}

	return nil
}

func (c *Client[T]) Delete(index []string) error {
	if _, err := c.client.Indices.Delete(index); err != nil {
		return err
	}

	return nil
}

func (c *Client[T]) Create(
	ctx context.Context,
	index string,
	data []byte,
) error {
	res, err := c.client.Indices.Create(
		index,
		c.client.Indices.Create.WithBody(bytes.NewReader(data)),
		c.client.Indices.Create.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("creating index: %s", res.String())
	}

	return nil
}

func (c *Client[T]) Bulk(
	ctx context.Context,
	index string,
	data []Data[T],
) error {
	dataBytes, err := Bytes(data)
	if err != nil {
		return err
	}

	res, err := c.client.Bulk(
		bytes.NewReader(dataBytes),
		c.client.Bulk.WithContext(ctx),
		c.client.Bulk.WithIndex(index),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk indexing: %s", res.String())
	}

	return nil
}

func (c *Client[T]) Count(index string) (int, error) {
	res, err := c.client.Count(
		c.client.Count.WithIndex(index),
	)
	if err != nil {
		return -1, err
	}
	defer res.Body.Close()

	type Response struct {
		Count  int `json:"count"`
		Shards struct {
			Total      int `json:"total"`
			Successful int `json:"successful"`
			Skipped    int `json:"skipped"`
			Failed     int `json:"failed"`
		} `json:"_shards"`
	}

	var resp Response
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return -1, err
	}

	return resp.Count, nil
}

func (c *Client[T]) CatIndex() ([]CatIndex, error) {
	res, err := c.client.Cat.Indices(
		c.client.Cat.Indices.WithFormat("json"),
		c.client.Cat.Indices.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var indices []CatIndex
	if err := json.NewDecoder(res.Body).Decode(&indices); err != nil {
		return nil, err
	}

	return indices, nil
}

func (c *Client[T]) Search(
	ctx context.Context,
	index string,
	query Query,
) (*SearchResult[T], error) {
	queryBytes, err := query.Bytes()
	if err != nil {
		return nil, err
	}

	res, err := c.client.Search(
		c.client.Search.WithContext(ctx),
		c.client.Search.WithIndex(index),
		c.client.Search.WithBody(bytes.NewReader(queryBytes)),
		c.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("searching documents: %s", res.String())
	}

	var result SearchResult[T]
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
