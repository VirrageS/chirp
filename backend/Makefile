.PHONY: all install test format

all: install test

install:
	go get -t -v ./...
	go install

test:
	go get -t -v ./...
	docker-compose -f ../docker/test.yml build && docker-compose -f ../docker/test.yml up -d
	until (pg_isready -h localhost -p 5433); do sleep 0.1; done
	go test ./...
	docker-compose -f ../docker/test.yml down -v

format:
	gofmt -w -s .
	go vet .
	golint ./...
