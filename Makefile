LDFLAGS ?=-s -w -X main.appVersion=dev-$(shell git rev-parse --short HEAD)-$(shell date +%y-%m-%d)
OUT ?= ./build
PROJECT ?=$(shell basename $(PWD))
SRC ?= ./cmd/$(PROJECT)
BINARY ?= $(OUT)/$(PROJECT)
export BUNDLE_TAG ?= latest
export PREFIX ?= manual

all: build lint test

build:
	mkdir -p $(OUT)
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -trimpath -o $(BINARY) $(SRC)

$(BINARY): build

run:
	CGO_ENABLED=0 go run -ldflags "$(LDFLAGS)" -trimpath $(SRC)

lint:
	golangci-lint run

test:
	go test ./...

.PHONY: all build run lint test

$(V).SILENT:
