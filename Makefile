SHELL := /bin/bash

build-docker-cli:
	go build -o bin/docker-cli ./cmd/docker-cli

install-docker-cli:
	go install ./cmd/docker-cli

test:
	go test ./...