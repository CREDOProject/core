VERSION := $(shell if [ -n "$$(git status --porcelain)" ]; then \
               echo "unreleased-$(shell git rev-parse --short HEAD)"; \
            elif git describe --tags --exact-match >/dev/null 2>&1; then \
               git describe --tags; \
            else \
               git rev-parse --short HEAD; \
            fi)
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
