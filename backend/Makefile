.PHONY: all install test format

all: install test

install:
	@echo "Getting dependencies..."
	@go get -t -v ./...
	@echo "Installing..."
	@go install

test:
	@echo "Starting docker..."
	@docker-compose -f ../docker/test.yml up --build -d &> /dev/null
	@echo "Waiting for docker instances..."
	@until (pg_isready -h localhost -p 5433 > /dev/null); do sleep 0.1; done
	@echo "Testing..."
	@CHIRP_CONFIG=test go test ./...; result=$$?; echo "Cleaning up..."; docker-compose -f ../docker/test.yml down -v &> /dev/null; exit $$result

format:
	@echo "Formatting..."
	@gofmt -w -s .
	@echo "Vetting..."
	@go vet .
	@echo "Linting..."
	@golint ./...
