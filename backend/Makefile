install:
	go install

test:
	docker-compose -f ../docker/test.yml build && docker-compose -f ../docker/test.yml up -d
	go test ./...
	docker-compose -f ../docker/test.yml down -v