t="coverprofile.txt"

default: build

clean:
	rm -f credo
	go mod download
	
build: clean
	go build -o credo

test:
	go test ./... -cover

t="coverprofile.txt"
coverage:
	go test -coverprofile=$t ./... && go tool cover -html=$t && unlink $t

debug: clean
	go build -gcflags "-N -l" -o credo
