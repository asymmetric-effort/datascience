SHELL := /bin/bash
.PHONY: all build test lint fmt clean website website-dev

SKILLS_VERSION := v0.0.28

all: lint test build

build:
	go build -o bin/pgmgo ./cmd/pgmgo

test:
	go test ./...

lint:
	go vet ./...
	gofmt -l .

fmt:
	gofmt -w .

clean:
	rm -rf bin/ website/dist/

website:
	cd website && node node_modules/.bin/vite build

website-dev:
	cd website && node node_modules/.bin/vite

install-skills:
	curl -fsSL https://skills.asymmetric-effort.com/install.sh | sh -s $(SKILLS_VERSION)
