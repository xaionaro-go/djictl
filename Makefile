
all: djictl-linux-amd64 djictl-linux-arm64

djictl-linux-amd64: builddir
	GOOS=linux GOARCH=amd64 go build -o build/djictl-linux-amd64 ./cmd/djictl/

djictl-linux-arm64: builddir
	GOOS=linux GOARCH=arm64 go build -o build/djictl-linux-arm64 ./cmd/djictl/

builddir:
	mkdir -p build
