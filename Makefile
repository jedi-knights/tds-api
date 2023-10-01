.PHONY: clean docs deps build mocks test lint

clean:
	rm -f tds-api
	rm -f junit.xml
	rm -rf zbuild

docs:
	~/go/bin/swag init -g main.go

deps:
	go mod tidy
	go mod download
	go install github.com/onsi/ginkgo/v2/ginkgo
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

build: clean
	go build main.go
	ls -l tds-api

build-prod: clean
	go build -ldflags="-s -w" main.go
	ls -l tds-api

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o zbuild/linux/386/${PWD##*/} -v ./...
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o zbuild/linux/amd64/${PWD##*/} -v ./...
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o zbuild/linux/arm64/${PWD##*/} -v ./...

build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o zbuild/darwin/amd64/${PWD##*/} -v ./...
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o zbuild/darwin/arm64/${PWD##*/} -v ./...

build-docker-all: clean build-linux build-darwin

build-docker:
	docker build -t tds-api .

docker-run:
	docker run -i -t -p 3000:8080 tds-api

docker-prune:
	docker image prune -a -f


mocks:
	go generate ./...

test: clean mocks
	ginkgo --junit-report=junit.xml ./...

lint:
	golangci-lint run ./...

run: docs
	go run main.go