APP=auth-api

.PHONY: run build test fmt tidy clean docker

run:
	go run ./cmd/server

build:
	go build -o bin/$(APP) ./cmd/server

test:
	go test ./... -v

fmt:
	go fmt ./...

tidy:
	go mod tidy

vendor:
	go mod vendor

deps:
	go mod tidy 
	go mod vendor

clean:
	rm -rf bin

docker:
	docker build -t $(APP) .
