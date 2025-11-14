package osr

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

type Client struct {
	osc *opensearch.Client
}

func NewClient(address []string, username, password string) (*Client, error) {
	client, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Addresses: address,
		Username:  username,
		Password:  password,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		osc: client,
	}, nil
}

func (c *Client) Ping() error {
	if _, err := c.osc.Ping(); err != nil {
		return err
	}

	return nil
}

func (c *Client) CatPlugins(ctx context.Context) ([]Plugin, error) {
	resp, err := opensearchapi.CatPluginsRequest{
		Format: "json",
	}.Do(ctx, c.osc)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return nil, fmt.Errorf("cat plugins: %s", resp.String())
	}

	var result []Plugin
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode: %v", err)
	}

	return result, nil
}

func (c *Client) CatIndex(ctx context.Context) ([]CatIndex, error) {
	resp, err := opensearchapi.CatIndicesRequest{
		Format: "json",
	}.Do(ctx, c.osc)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return nil, fmt.Errorf("cat indices: %s", resp.String())
	}

	var result []CatIndex
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode: %v", err)
	}

	return result, nil
}

func (c *Client) Count(
	ctx context.Context,
	indexName string,
) (int, error) {
	resp, err := opensearchapi.CountRequest{
		Index: []string{indexName},
	}.Do(ctx, c.osc)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return -1, fmt.Errorf("count: %s", resp.String())
	}

	type Result struct {
		Count int `json:"count"`
	}

	var result Result
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return -1, fmt.Errorf("decode: %v", err)
	}

	return result.Count, nil
}

func (c *Client) Delete(
	ctx context.Context,
	indexName []string,
) error {
	resp, err := opensearchapi.IndicesDeleteRequest{
		Index: indexName,
	}.Do(ctx, c.osc)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return fmt.Errorf("delete index: %s", resp.String())
	}

	return nil
}

func (c *Client) Create(
	ctx context.Context,
	indexName string,
	body []byte,
) error {
	resp, err := opensearchapi.IndicesCreateRequest{
		Index: indexName,
		Body:  bytes.NewReader(body),
	}.Do(ctx, c.osc)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return fmt.Errorf("create index: %s", resp.String())
	}

	return nil
}

func (c *Client) Refresh(
	ctx context.Context,
	indexName []string,
) error {
	resp, err := opensearchapi.IndicesRefreshRequest{
		Index: indexName,
	}.Do(ctx, c.osc)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return fmt.Errorf("refresh: %s", resp.String())
	}

	return nil
}

func Bulk[T any](
	ctx context.Context,
	client *Client,
	data []byte,
) error {
	if _, err := Decode[T](data); err != nil {
		return err
	}

	resp, err := client.osc.Bulk(
		bytes.NewReader(data),
		client.osc.Bulk.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return fmt.Errorf("bulk indexing: %s", resp.String())
	}

	return nil
}

func Search[T any](
	ctx context.Context,
	client *Client,
	indexName []string,
	query []byte,
) (*SearchResult[T], error) {
	resp, err := opensearchapi.SearchRequest{
		Index: indexName,
		Body:  bytes.NewReader(query),
	}.Do(ctx, client.osc)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return nil, fmt.Errorf("search: %s", resp.String())
	}

	var result SearchResult[T]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode: %v", err)
	}

	return &result, nil
}
