package osr

import (
	"context"
	"crypto/tls"
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

func (c *Client) Delete(ctx context.Context, index []string) error {
	resp, err := opensearchapi.IndicesDeleteRequest{
		Index: index,
	}.Do(ctx, c.osc)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) Create(ctx context.Context, index string, body io.Reader) error {
	resp, err := opensearchapi.IndicesCreateRequest{
		Index: index,
		Body:  body,
	}.Do(ctx, c.osc)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
