VERSION=0.5.0
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_HASH=$(shell git rev-parse HEAD)

RELEASES=bin/xliffer-$(VERSION).linux.amd64 \
		 bin/xliffer-$(VERSION).windows.amd64.exe \
		 bin/xliffer-$(VERSION).freebsd.amd64 \
		 bin/xliffer-$(VERSION).darwin.amd64

LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildDate=$(BUILD_DATE) -X main.GitHash=$(GIT_HASH)"

# something to support old muscle memory: "make" :)
build: force
	go build -v

releases: $(RELEASES)

bin/xliffer-$(VERSION).linux.amd64: bin
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $@

bin/xliffer-$(VERSION).windows.amd64.exe: bin
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $@

bin/xliffer-$(VERSION).darwin.amd64: bin
	env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $@

bin/xliffer-$(VERSION).freebsd.amd64: bin
	env GOOS=freebsd GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $@

bin:
	mkdir $@

test:
	go test

.PHONY : force
