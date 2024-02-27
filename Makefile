default: build

clean:
	rm -f credo
	go mod download
	
build: clean
	go build -o credo

coverage:
	go test -coverprofile=$t ./... && go tool cover -html=$t && unlink $t


debug: clean
	go build -gcflags "-N -l" -o credo
