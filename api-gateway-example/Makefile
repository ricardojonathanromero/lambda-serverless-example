.PHONY: build clean deploy

dependencies:
	go get github.com/DataDog/datadog-lambda-go

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ./bin/create ./todo-create/cmd/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ./bin/delete ./todo-delete/cmd/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ./bin/get ./todo-get/cmd/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ./bin/list ./todo-list/cmd/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ./bin/update ./todo-update/cmd/main.go

clean:
	rm -rf ./bin

tests:
	go test ./...

deploy: clean build
	sls deploy --verbose
