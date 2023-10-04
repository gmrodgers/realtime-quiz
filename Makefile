.DEFAULT_GOAL := help

rabbitmq-cluster: kill-rabbitmq-cluster
	# TODO: does this need --privileged
	# TODO: why does this need on wsl
	docker run --privileged -d --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.12-management

kill-rabbitmq-cluster:
	docker rm rabbitmq --force

server: ## run server with 1 player
	go run cmd/server/main.go -n 1

client: ## run client with name glen
	go run cmd/client/main.go name-glen

test: rabbitmq-cluster ## spin up rmq, run server and client, spin down rmq
	sleep 15
	$(MAKE) server &
	$(MAKE) client
	$(MAKE) kill-rabbitmq-cluster

help:  ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'