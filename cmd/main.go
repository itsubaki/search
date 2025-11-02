package main

import (
	"fmt"
	"io"

	"github.com/elastic/go-elasticsearch/v9"
)

func main() {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		panic(err)
	}

	resp, err := client.Info()
	if err != nil && err != io.EOF {
		panic(err)
	}

	fmt.Println(resp)

	pong, err := client.Ping()
	if err != nil && err != io.EOF {
		panic(err)
	}

	fmt.Println(pong)
}
