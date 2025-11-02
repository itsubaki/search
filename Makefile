SHELL := /bin/bash

elastic:
	docker run --rm --name elasticsearch --net elastic -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" elasticsearch:9.2.0

init:
	docker network create elastic

clean:
	docker network remove elastic
