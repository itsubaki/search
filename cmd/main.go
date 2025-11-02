package main

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v9"
)

func main() {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		panic(err)
	}

	pong, err := client.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println(pong)
}
