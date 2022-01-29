BIN := "./bin/sysmon"
DOCKER_IMG="sysmon:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/sysmon

run: build
	$(BIN) -config ./configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

run-img-host: build-img
	docker run --network host $(DOCKER_IMG)


version: build
	$(BIN) version

test:
	go test -race ./...

test-all:
	go test --tags longtest -race ./...

test-all-clean-cache:
	go clean -testcache
	go test --tags longtest -race ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

lint-fix: install-lint-deps
	golangci-lint run ./... --fix

install-gen-deps:
	(which protoc-gen-go > /dev/null) || go install google.golang.org/protobuf/cmd/protoc-gen-go
	(which protoc-gen-go-grpc > /dev/null) || go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	(which protoc-gen-grpc-gateway > /dev/null) || go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	(which protoc-gen-openapiv2 > /dev/null) || go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2

generate: install-gen-deps
	go generate ./...

.PHONY: build run build-img run-img version test lint
