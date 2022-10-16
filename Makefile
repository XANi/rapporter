# generate version number
version=$(shell git describe --tags --long --always --dirty|sed 's/^v//')
binfile=rapporter

all:
	go build -ldflags "-X main.version=$(version)" $(binfile).go
	-@go fmt

static:
	CGO_ENABLED=0 go build -ldflags "-X main.version=$(version) -extldflags \"-static\"" -o $(binfile).static $(binfile).go

arch:
	mkdir -p bin
	CGO_ENABLED=0 GOARCH=arm go build  -ldflags "-X main.version=$(version) -extldflags \"-static\"" -o bin/$(binfile).arm $(binfile).go
	CGO_ENABLED=0 GOARCH=arm64 go build  -ldflags "-X main.version=$(version) -extldflags \"-static\"" -o bin/$(binfile).aarch64 $(binfile).go
	CGO_ENABLED=0 GOARCH=amd64 go build  -ldflags "-X main.version=$(version) -extldflags \"-static\"" -o bin/$(binfile).amd64 $(binfile).go
	CGO_ENABLED=0 GOARCH=386 go build  -ldflags "-X main.version=$(version) -extldflags \"-static\"" -o bin/$(binfile).386 $(binfile).go
	sha256sum bin/dpp.* > bin/Checksum
