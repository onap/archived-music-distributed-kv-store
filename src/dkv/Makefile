GOPATH := $(shell realpath "$(PWD)/../../")
BINARY := dkv
DEPENDENCIES := github.com/golang/dep/cmd/dep

export GOPATH ...

all: test build
deploy: test build

build: deps format
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w' -o $(GOPATH)/target/$(BINARY) -v main.go

clean:
	go clean
	rm -f $(GOPATH)/target/$(BINARY)

test: deps
	go test -cover ./...

format:
	go fmt ./...

deps:
	go get -u $(DEPENDENCIES)
	$(GOPATH)/bin/dep ensure

.PHONY: test
