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

update-usage: build
	head -n $(shell grep -nE '^## Usage' README.md | tr ':' ' ' | awk '{print $$1}') README.md > README.md.temp
	tail -n +$(shell grep -nE '^## Contributing' README.md | tr ':' ' ' | awk '{print $$1}') README.md > README.md.tail
	echo >> README.md.temp
	echo '```' >> README.md.temp
	$(BINARY) --help  >> README.md.temp || true
	echo '```' >> README.md.temp
	mv -f README.md.temp README.md
	echo >> README.md
	cat README.md.tail >> README.md
	rm README.md.tail

.PHONY: all build run lint test update-usage

$(V).SILENT:
