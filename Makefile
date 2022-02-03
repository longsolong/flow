PROJECT_NAME := $(APP_NAME)
MAIN_PACKAGE := "cmd/apiserver"
PKG_LIST := $(shell go list ./... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)
WIRE_FILES := $(shell find . -name 'wire.go' | grep -v /vendor/)

# https://github.com/golangci/awesome-go-linters
LINTERS := \
	golang.org/x/lint/golint \
	honnef.co/go/tools/cmd/staticcheck

.PHONY: all init dep build clean test coverage coverhtml lint golint vet staticcheck

all: build

init: dep testdep ## Download dependencies and add git hooks
	find .git/hooks -type l -exec rm {} \;
	find githooks -type f -exec ln -sf ../../{} .git/hooks/ \;
	go env -w GOPROXY=https://goproxy.cn,direct

lint: testdep ## Lint files
	@`go env GOPATH`/bin/golint -set_exit_status ${PKG_LIST}

vet: testdep ## Checks correctness 
	@go vet ${PKG_LIST}

staticcheck: testdep ## Analyses code
	@`go env GOPATH`/bin/staticcheck ${PKG_LIST}

install-golangci-lint: ## 在项目目录外手动执行
	@go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.35.2

golangci-lint: testdep
	@`go env GOPATH`/bin/golangci-lint run

test: ## Run unit tests
	@go test -short ${PKG_LIST}

test-coverage: ## Run tests with coverage
	@go test -short -coverprofile cover.out -covermode=atomic ${PKG_LIST} 
	@cat cover.out >> coverage.txt

test-int: ## Run unit and integration tests
	@go test -short -tags=integration ${PKG_LIST}

coverage: ## Generate global code coverage report
	./scripts/coverage.sh;

coverhtml: ## Generate global code coverage report in HTML
	./scripts/coverage.sh html;

dep: ## Get dependencies
	@go mod tidy
	@go mod vendor

testdep: ## Get dev dependencies
	@go get -v $(LINTERS)

run:
	./bin/$(PROJECT_NAME) serve

build: dep ## Build the binary file
	@go build -i -o ./bin/$(PROJECT_NAME) ./$(MAIN_PACKAGE)

wiregen:
	echo $(WIRE_FILES) | xargs sh -c 'for filename; do `go env GOPATH`/bin/wire "$$(dirname $$filename)"; done' sh

generate:
	@go generate

dev:
	make build
	make run

clean: ## Remove previous build
	@rm -f ./bin

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
