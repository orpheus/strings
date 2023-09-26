.PHONY: help build start printos

timestamp := $(shell date +'%Y_%m_%d_%H_%M_%S')

help: ## : Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_%-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' ${MAKEFILE_LIST}

build: ## : Build dependencies
	go mod tidy
	GOOS=darwin GOARCH=amd64 go build -o build/strings main.go

start: build ## : Start the client
	DB_USER=postgres \
    DB_PASS= \
    DB_NAME=strings \
    DB_HOST=localhost \
    DB_PORT=5432 \
    ./build/strings

run: ## : Run app without build
	DB_USER=postgres \
	DB_PASS= \
	DB_NAME=strings \
	DB_HOST=localhost \
	DB_PORT=5432 \
	./build/strings

dump: ## : dump postgres database
	/usr/local/bin/pg_dump --dbname=strings --file="${HOME}/strings_localhost-$(timestamp)-dump.sql" --username=postgres --host=localhost --port=5432


docker: ## : Make docker image
	rm -f build/strings-docker
	go mod tidy
	# https://stackoverflow.com/questions/34729748/installed-go-binary-not-found-in-path-on-alpine-linux-docker
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on go build -o build/strings-linux-amd64 main.go
	docker build -t psytek/strings:local .

docker-run: ## : Run docker image
	docker run -it \
		-p 8081:8080 \
		-e DB_HOST=host.docker.internal \
		-e DB_PORT=5432 \
		-e DB_NAME=strings \
		-e DB_USER=postgres \
		-e DB_PASS=password \
		psytek/strings:local

docker-exec: ## : Run docker image
	docker run -it psytek/strings:local sh