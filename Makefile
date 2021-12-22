BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: all help

gen-proto: ## Generate protobuf classes
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/adder/v1/adder.proto

server-image: ## Build the server image
	docker build -t bradbeck/add-service -f service/Dockerfile .

api-image: ## Build the api image
	docker build -t bradbeck/api-service -f api/Dockerfile .

all-images: server-image api-image ## Build both server and api images

all: help
help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
