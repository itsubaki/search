SHELL := /bin/bash

run:
	go run cmd/main.go

elastic:
	docker run --rm --name elasticsearch -p 9200:9200 -e "discovery.type=single-node" -e "ES_JAVA_OPTS=-Xms1g -Xmx1g" elasticsearch:9.2.0

init:
	docker exec -it elasticsearch bin/elasticsearch-plugin install analysis-icu
	docker exec -it elasticsearch bin/elasticsearch-plugin install analysis-kuromoji
	docker restart elasticsearch
	docker exec -it elasticsearch bin/elasticsearch-reset-password -u elastic
