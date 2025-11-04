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

type Client struct {
	es *elasticsearch.Client
}

func NewClient(address []string, username, password string) (*Client, error) {
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

	return &Client{
		es: es,
	}, nil
}

func (c *Client) Ping() error {
	if _, err := c.es.Ping(); err != nil {
		return err
	}

	return nil
}

func (c *Client) Delete(indexNames []string) error {
	if _, err := c.es.Indices.Delete(indexNames); err != nil {
		return err
	}

	return nil
}

func (c *Client) Create(
	ctx context.Context,
	indexName string,
	index Index,
) error {
	data, err := json.Marshal(index)
	if err != nil {
		return err
	}

	res, err := c.es.Indices.Create(
		indexName,
		c.es.Indices.Create.WithBody(bytes.NewReader(data)),
		c.es.Indices.Create.WithContext(ctx),
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

func (c *Client) Count(indexName string) (int, error) {
	res, err := c.es.Count(
		c.es.Count.WithIndex(indexName),
	)
	if err != nil {
		return -1, err
	}
	defer res.Body.Close()

	type Response struct {
		Count int `json:"count"`
	}

	var resp Response
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return -1, err
	}

	return resp.Count, nil
}

func (c *Client) CatIndex() ([]CatIndex, error) {
	res, err := c.es.Cat.Indices(
		c.es.Cat.Indices.WithFormat("json"),
		c.es.Cat.Indices.WithPretty(),
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

func Bulk[T any](
	ctx context.Context,
	client *Client,
	indexName string,
	data []Data[T],
) error {
	dataBytes, err := Bytes(data)
	if err != nil {
		return err
	}

	res, err := client.es.Bulk(
		bytes.NewReader(dataBytes),
		client.es.Bulk.WithContext(ctx),
		client.es.Bulk.WithIndex(indexName),
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

func Search[T any](
	ctx context.Context,
	client *Client,
	indexName string,
	query Query,
) (*SearchResult[T], error) {
	queryBytes, err := query.Bytes()
	if err != nil {
		return nil, err
	}

	res, err := client.es.Search(
		client.es.Search.WithContext(ctx),
		client.es.Search.WithIndex(indexName),
		client.es.Search.WithBody(bytes.NewReader(queryBytes)),
		client.es.Search.WithTrackTotalHits(true),
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
