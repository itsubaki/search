SHELL := /bin/bash

elastic:
	docker run --rm --name elasticsearch --net elastic -p 9200:9200 -e "discovery.type=single-node" -e "ES_JAVA_OPTS=-Xms1g -Xmx1g" elasticsearch:9.2.0

init:
	docker network create elastic

clean:
	docker network remove elastic
