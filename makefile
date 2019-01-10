all: test build
build:
	go build -o etest-server app/main.go
	go build -o etest-cli cli/main.go
test: 
	go test -v ./...