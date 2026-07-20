SHELL = /bin/bash

MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
ROOT := $(dir $(MKFILE_PATH))
GOBIN ?= $(ROOT)/tools/bin
ENV_PATH = PATH=$(GOBIN):$(PATH)
BIN_PATH ?= $(ROOT)/bin

GOPRIVATE = GOPRIVATE=github.com/Mikhalevich/

LINTER_NAME := golangci-lint
LINTER_VERSION := v2.12.2

.PHONY: all build test compose-up compose-down load-test-data vendor install-linter \
lint fmt tools-update generate load-assets

all: build

build:
	go build -mod=vendor -o $(BIN_PATH)/bot ./cmd/bot/main.go
	go build -mod=vendor -o $(BIN_PATH)/poller ./cmd/poller/main.go

test:
	go test ./...

compose-up:
	docker compose -f ./script/docker/docker-compose.yml up --build

compose-down:
	docker compose -f ./script/docker/docker-compose.yml down

vendor:
	$(GOPRIVATE) go mod tidy
	$(GOPRIVATE) go mod vendor

install-linter:
	if [ ! -f $(GOBIN)/$(LINTER_VERSION)/$(LINTER_NAME) ]; then \
		echo INSTALLING $(GOBIN)/$(LINTER_VERSION)/$(LINTER_NAME) $(LINTER_VERSION) ; \
		curl -sSfL https://golangci-lint.run/install.sh  | sh -s -- -b $(GOBIN)/$(LINTER_VERSION) $(LINTER_VERSION) ; \
		echo DONE ; \
	fi

lint: install-linter
	$(GOBIN)/$(LINTER_VERSION)/$(LINTER_NAME) run --config .golangci.yml

fmt: install-linter
	$(GOBIN)/$(LINTER_VERSION)/$(LINTER_NAME) fmt --config .golangci.yml

tools-update:
	go get tool

generate:
	$(ENV_PATH) go generate ./...

load-assets:
	docker run -it --rm --network host \
		-v ./script/db/dataset/assets.sql:/script/assets.sql \
		alpine/psql:17.7 \
		"postgresql://bot:bot@localhost:5432/bot" -f /script/assets.sql
