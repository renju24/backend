vendor: ## download vendors
	go mod vendor

build: ## build api-server
	go build -mod=vendor -o ./build/apiserver -v ./cmd/apiserver

test: ## run tests
	go test ./...

.PHONY: build
.PHONY: vendor

help:
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
