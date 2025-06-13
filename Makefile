export RELEASE_VERSION ?= $(shell git describe --tags --always)
export GIT_COMMIT := $(shell git rev-parse HEAD)
export BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

t="coverprofile.txt"

default: build

clean:
	rm -f credo
	go mod download
	

.PHONY: build
build: clean
	go build -ldflags "\
		-X 'credo/version.Version=${RELEASE_VERSION}' \
		-X 'credo/version.Commit=${GIT_COMMIT}' \
		-X 'credo/version.BuildDate=${BUILD_DATE}'" \
		-o credo
test:
	go test ./... -cover

t="coverprofile.txt"
coverage:
	go test -coverprofile=$t ./... && go tool cover -html=$t && unlink $t

debug: clean
	go build -gcflags "-N -l" -o credo
