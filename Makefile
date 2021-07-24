APP := fastly
VERSION ?= $(shell cat VERSION)

all: vet fmt test build

verify: vet fmt test

vet:
	@go vet ./...

fmt:
	@go fmt ./...

test:
	@go test ./...

build:
	@docker build -t $(APP):$(VERSION) .
