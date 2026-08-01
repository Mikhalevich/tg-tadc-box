SHELL = /bin/bash

MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
ROOT := $(dir $(MKFILE_PATH))
GOBIN ?= $(ROOT)/tools/bin
ENV_PATH = PATH=$(GOBIN):$(PATH)
BIN_PATH ?= $(ROOT)/bin

GOPRIVATE = GOPRIVATE=github.com/Mikhalevich/

LINTER_NAME := golangci-lint
LINTER_VERSION := v2.12.2

APP_TAG := 0.1.2

.PHONY: all build test bench compose-up compose-down load-test-data vendor install-linter \
lint fmt tools-update generate load-assets \
minikube-load-images minikube-decrypt-secrets minikube-encrypt-secrets minikube-helm-install minikube-helm-uninstall \
install-helm-secrets generate-helm-secrets \
do-helm-load-images do-helm-encrypt-secrets do-helm-install do-helm-uninstall

all: build

build:
	go build -mod=vendor -o $(BIN_PATH)/bot ./cmd/bot/main.go
	go build -mod=vendor -o $(BIN_PATH)/poller ./cmd/poller/main.go
	go build -mod=vendor -o $(BIN_PATH)/outboxpoller ./cmd/outboxpoller/main.go
	go build -mod=vendor -o $(BIN_PATH)/botwebhook ./cmd/botwebhook/main.go

test:
	go test ./...

bench:
	go test -bench=. -benchmem ./...

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

minikube-load-images:
	./script/k8s/minikube-helm/load_images.sh ${APP_TAG}

minikube-decrypt-secrets:
	./script/k8s/minikube-helm/decrypt_secrets.sh

minikube-encrypt-secrets:
	./script/k8s/minikube-helm/encrypt_secrets.sh

minikube-helm-install:
	./script/k8s/minikube-helm/install.sh ${APP_TAG}

minikube-helm-uninstall:
	./script/k8s/minikube-helm/uninstall.sh

install-helm-secrets:
	helm plugin install https://github.com/jkroepke/helm-secrets/releases/download/v4.7.4/secrets-4.7.4.tgz --verify=false

generate-helm-secrets:
	mkdir -p ~/.config/sops/age/
	age-keygen -o ~/.config/sops/age/keys.txt

do-helm-load-images:
	./script/k8s/do-helm/load_images.sh ${APP_TAG}

do-helm-decrypt-secrets:
	./script/k8s/do-helm/decrypt_secrets.sh

do-helm-encrypt-secrets:
	./script/k8s/do-helm/encrypt_secrets.sh

do-helm-install:
	./script/k8s/do-helm/install.sh ${APP_TAG}

do-helm-uninstall:
	./script/k8s/do-helm/uninstall.sh
