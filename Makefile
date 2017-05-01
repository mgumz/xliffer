all: force
	go build -v

# create a binary with debug-symbols removed
all-release: force 
	go build -v -ldflags "-s"

xliffer.amd64.windows.exe:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 \
	go build -v -ldflags "-s" -a \
		-o $@ \
		.


test:
	go test

.PHONY : force
