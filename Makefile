.PHONY: help build start printos

help: ## : Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_%-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' ${MAKEFILE_LIST}

build: ## : Build dependencies
	go mod tidy
	GOOS=darwin GOARCH=amd64 go build -o build/strings main.go
	rm -f "${GOPATH}/bin/strings" && cp build/strings "${GOPATH}/bin"

start: build ## : Start the client
	DB_USER=postgres \
    DB_PASS= \
    DB_NAME=strings \
    DB_HOST=localhost \
    DB_PORT=5432 \
    strings


