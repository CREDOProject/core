VERSION := $(shell git describe --tags --exact-match 2>/dev/null || echo "unreleased-$(shell git rev-parse --short HEAD)")
COMMIT := $(shell git rev-parse HEAD)
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

t="coverprofile.txt"

default: build

clean:
	rm -f credo
	go mod download
	

.PHONY: build
build: clean
	go build -ldflags "\
		-X 'credo/version.Version=${VERSION}' \
		-X 'credo/version.Commit=${COMMIT}' \
		-X 'credo/version.BuildDate=${BUILD_DATE}'" \
		-o credo
test:
	go test ./... -cover

t="coverprofile.txt"
coverage:
	go test -coverprofile=$t ./... && go tool cover -html=$t && unlink $t

debug: clean
	go build -gcflags "-N -l" -o credo
