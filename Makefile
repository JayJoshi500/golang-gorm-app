APP=auth-api

.PHONY: run build test fmt tidy clean docker docker-down

run:
	go run ./cmd/api

build:
	go build -o bin/$(APP) ./cmd/api

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
	docker compose up --build

docker-down:
	docker compose down
