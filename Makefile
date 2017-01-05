VERSION := $(shell git describe --tags --always --dirty)

.PHONY: bin bin/fetch

all: bin/fetch

bin:
	mkdir -p bin

bin/fetch: bin
		go build -i -v -o $@ -ldflags "-X main.Version=$(VERSION)" ./cmd/fetch
