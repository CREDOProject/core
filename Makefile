build:
	rm -f credo
	go mod download
	go build -o credo

