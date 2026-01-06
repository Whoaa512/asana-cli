.PHONY: build test lint clean install

BINARY := asana
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o $(BINARY) ./cmd/asana

test:
	go test -race -v ./...

lint:
	golangci-lint run

clean:
	rm -f $(BINARY)
	go clean -testcache

install:
	go install $(LDFLAGS) ./cmd/asana

# Run all checks (used by CI)
check: lint test
