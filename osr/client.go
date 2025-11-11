package osr

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
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

func (c *Client) CatIndex() ([]CatIndex, error) {
	resp, err := opensearchapi.CatIndicesRequest{
		Format: "json",
	}.Do(context.Background(), c.osc)
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
	index string,
) (int, error) {
	resp, err := opensearchapi.CountRequest{
		Index: []string{index},
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
	index []string,
) error {
	resp, err := opensearchapi.IndicesDeleteRequest{
		Index: index,
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
	index string,
	body io.Reader,
) error {
	resp, err := opensearchapi.IndicesCreateRequest{
		Index: index,
		Body:  body,
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
	index []string,
) error {
	resp, err := opensearchapi.IndicesRefreshRequest{
		Index: index,
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
	data []Data[T],
) error {
	b, err := Bytes(data)
	if err != nil {
		return err
	}

	resp, err := client.osc.Bulk(
		bytes.NewReader(b),
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
	index []string,
	query io.Reader,
) (*SearchResult[T], error) {
	resp, err := opensearchapi.SearchRequest{
		Index: index,
		Body:  query,
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

func MustRead(r io.Reader) string {
	bytes, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}
