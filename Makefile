ifneq (,$(wildcard .env))
    include .env
    export
endif


IMAGE_NAME := customer-service
IMAGE_TAG := latest

format:
	gofmt -w .

# brew install golangci-lint
linter:
	golangci-lint run ./...
	#TODO: fix the linter errors

integration: # TODO Higher tests coverage
	go test -v -coverprofile cover.out ./... && \
    go tool cover -html=cover.out -o cover.html && \
    open cover.html

integration-ci:
	go test -v ./...

unit:
	go test -v ./pkg/server

# Should be ran every time customer interface changes
regenerate-mocks:
	mockery --name=CustomerRepository --dir=./pkg/repository --output=./mocks --outpkg=mocks

# Create the kind cluster
create-cluster:
	kind create cluster --name customer-service-cluster

# Build the Docker image
build-image:
	docker build -t customer-service:latest .

# Load the Docker image into kind cluster
load-image: build-image
	kind load docker-image customer-service:latest --name customer-service-cluster

# Deploy the Helm chart
deploy: load-image
	helm upgrade --install -f ./helm/values.yaml customer-service ./helm

# Clean up the kind cluster
delete-cluster:
	kind delete cluster --name customer-service-cluster
