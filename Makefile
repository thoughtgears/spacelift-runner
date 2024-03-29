#!make
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

OS = $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH = $(shell uname -m)
GIT_COMMIT = $(shell git rev-parse --short HEAD)
GIT_REPO = "thoughtgears/spacelift-runner"
SERVICE_NAME = huston

.PHONY: clean build test docker

all: clean build

clean:
	@rm -rf builds

build:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -a -installsuffix cgo -ldflags="-X main.Version=$(GIT_COMMIT)" -o builds/$(SERVICE_NAME)-$(OS)-$(ARCH)

build-ci: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-X main.Version=$(GIT_COMMIT)" -o builds/$(SERVICE_NAME)-linux-amd64
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-X main.Version=$(GIT_COMMIT)" -o builds/$(SERVICE_NAME)-darwin-amd64
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -a -installsuffix cgo -ldflags="-X main.Version=$(GIT_COMMIT)" -o builds/$(SERVICE_NAME)-darwin-arm64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-X main.Version=$(GIT_COMMIT)" -o builds/$(SERVICE_NAME)-windows-amd64.exe

build-docker: clean build-ci
	docker build -t ghcr.io/$(GIT_REPO):$(GIT_COMMIT) .
