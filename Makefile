include .env
export

IMAGE_NAME := customer-service
IMAGE_TAG := latest

format:
	gofmt -w .

linter:
	golint ./...

integration: # TODO Higher tests coverage
	go test -v -coverprofile cover.out ./... && \
    go tool cover -html=cover.out -o cover.html && \
    open cover.html

unit:
	go test -v ./pkg/server

# Should be ran every time customer interface changes
regenerate-mocks:
	mockery --name=CustomerRepository --dir=./pkg/repository --output=./mocks --outpkg=mocks
