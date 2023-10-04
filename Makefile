.DEFAULT_GOAL := help

rabbitmq-cluster:
	# TODO: does this need --privileged
	# TODO: why does this need on wsl
	docker run --privileged -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.12-management

kill-rabbitmq-cluster:
	docker rm rabbitmq

server: players=3
server: ## run server with number of players players=<insert>
	go run cmd/server/main.go -n $(players)

client: name=glen
client: ## run client with name name=<insert>
	go run cmd/client/main.go -n $(name)

help:  ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'