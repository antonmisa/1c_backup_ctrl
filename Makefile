include .env
export

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

test: ### run test
	go test -v -cover -race ./internal/...
.PHONY: test

mock: ### run mockgen
	go generate ./...
.PHONY: mock

build: ### build for windows, linux & darwin GOOS, all is x64  
	env GOOS="linux" GOARCH="amd64" CGO_ENABLED=0 go build -o build/ctrl_linux -ldflags "-w -s" cmd/app/main.go
	env GOOS="darwin" GOARCH="amd64" CGO_ENABLED=0 go build -o build/ctrl_darwin -ldflags "-w -s" cmd/app/main.go
	env GOOS="windows" GOARCH="amd64" CGO_ENABLED=0 go build -o build/ctrl_win64 -ldflags "-w -s" cmd/app/main.go
.PHONY: build