SHELL := /bin/bash

OS_NAME := $(shell uname -s | tr A-Z a-z)

build-docker-cli:
	go build -o bin/docker-cli ./cmd/docker-cli

install-docker-cli:
	go install ./cmd/docker-cli
	@if [ $(OS_NAME) = "linux" ]; then \
		docker-cli completion bash > /tmp/completion; \
	else \
		docker-cli completion zsh > /tmp/completion; \
	fi

test:
	go test ./...
