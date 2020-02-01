NAME=und
BINDIR=bin
LDFLAGS=-s -w -X "main.version=$(shell git describe --tags --dirty --always)"
GOBUILD=CGO_ENABLED=0 go build --ldflags="$(LDFLAGS)" -trimpath

all: linux-arm64 linux-amd64 darwin-amd64 windows-amd64

linux-arm64:
	GOARCH=arm64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-amd64:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

darwin-amd64:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

windows-amd64:
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe

check: errcheck golint
	go mod tidy
	gofmt -s -l -w .
	errcheck -blank
	golint
	go vet

errcheck:
	go get -u github.com/kisielk/errcheck

golint:
	go get -u golang.org/x/lint/golint
