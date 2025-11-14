SHELL := /bin/bash

run:
	go run cmd/main.go

up:
	docker run --rm --name opensearch -p 9200:9200 -p 9600:9600 -e "discovery.type=single-node" -e "OPENSEARCH_JAVA_OPTS=-Xms1g -Xmx1g" -e "OPENSEARCH_INITIAL_ADMIN_PASSWORD=xuYz3_cAXYh7" opensearchproject/opensearch:latest

init:
	docker exec -it opensearch bin/opensearch-plugin list
	docker exec -it opensearch bin/opensearch-plugin install analysis-icu
	docker exec -it opensearch bin/opensearch-plugin install analysis-kuromoji
	docker restart opensearch

log:
	docker logs -f opensearch
