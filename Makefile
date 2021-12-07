build:
	go mod download
	gofmt -l -s -w .
	go test -v ./...
	go mod tidy

update-deps:
	go get -u ./...
	go mod tidy