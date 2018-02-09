GOPATH := $(shell realpath "$(PWD)/../../")
BINARY=dkv
DEPENDENCIES := github.com/golang/dep/cmd/dep

export GOPATH ...

all: build test
deploy: build test

build: deps format
	echo $(GOPATH)
	$(GOPATH)/bin/dep ensure
	go build -o $(GOPATH)/target/$(BINARY) -v main.go

clean:
	go clean
	rm -f $(GOPATH)/target/$(BINARY)

test:
	go test -v ./api/...

format:
	go fmt ./api/...

deps:
	go get -u $(DEPENDENCIES)
