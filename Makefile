SHELL := /bin/bash

run:
	go run cmd/main.go

elastic:
	docker run --rm --name elasticsearch -p 9200:9200 -e "discovery.type=single-node" -e "ES_JAVA_OPTS=-Xms1g -Xmx1g" elasticsearch:9.2.0

elastic-init:
	docker exec -it elasticsearch bin/elasticsearch-plugin install analysis-icu
	docker exec -it elasticsearch bin/elasticsearch-plugin install analysis-kuromoji
	docker restart elasticsearch
	docker exec -it elasticsearch bin/elasticsearch-reset-password -u elastic


opensearch:
	docker run --rm --name opensearch -p 9200:9200 -p 9600:9600 -e "discovery.type=single-node" -e "OPENSEARCH_JAVA_OPTS=-Xms1g -Xmx1g" -e "OPENSEARCH_INITIAL_ADMIN_PASSWORD=xuYz3_cAXYh7" opensearchproject/opensearch:latest

opensearch-init:
	docker exec -it opensearch bin/opensearch-plugin install analysis-icu
	docker exec -it opensearch bin/opensearch-plugin install analysis-kuromoji
	docker restart opensearch

opensearch-log:
	docker logs -f opensearch
