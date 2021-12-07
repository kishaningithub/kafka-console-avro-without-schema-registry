build:
	go mod download
	gofmt -l -s -w .
	go test -v ./...
	go mod tidy