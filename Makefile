PROJECT := $(shell go list ./... | head -n 1)

all: lint vet test

.PHONY: help
help: ## print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test:
	go test $(PKG_LIST) -v -cover

.PHONY: coverage
coverage:
	go test -v -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	docker run \
		-it \
		--rm \
		-v "$(CURDIR):/go/src/$(PROJECT)" \
		-w "/go/src/$(PROJECT)" \
		golangci/golangci-lint:v1.19.1 \
		golangci-lint run
