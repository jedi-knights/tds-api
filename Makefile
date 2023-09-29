.PHONY: clean docs deps build mocks test lint

clean:
	rm -f tds
	rm -f junit.xml

docs:
	~/go/bin/swag init -g main.go

deps:
	go mod tidy
	go mod download
	go install github.com/onsi/ginkgo/v2/ginkgo
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

build: clean
	go build -v ./...

mocks:
	go generate ./...

test: clean mocks
	ginkgo --junit-report=junit.xml ./...

lint:
	golangci-lint run ./...

run: docs
	go run main.go