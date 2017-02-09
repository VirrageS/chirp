.PHONY: core frontend backend production clean

core:
	@echo "Starting core..."
	@docker-compose -f docker/core.yml up --build -d &> /dev/null
	@echo "Done" # TODO: actualy not done... core is not ready yet

frontend:
	@echo "Starting frontend..."
	@docker-compose -f docker/frontend.yml up --build -d &> /dev/null
	@echo "Done" # TODO: actualy not done... frontend is not ready yet

backend:
	@echo "Starting backend..."
	@docker-compose -f docker/backend.yml up --build -d &> /dev/null
	@echo "Done" # TODO: actualy not done... backend is not ready yet

production:
	@echo "Starting production..."
	@docker-compose -f docker/production.yml up --build -d &> /dev/null
	@echo "Done" # TODO: actualy not done... production is not ready yet

clean:
	@echo "Cleaning up..."
	@docker-compose -f docker/core.yml down -v &> /dev/null
	@docker-compose -f docker/test.yml down -v &> /dev/null
	@docker-compose -f docker/backend.yml down -v &> /dev/null
	@docker-compose -f docker/frontend.yml down -v &> /dev/null
	@docker-compose -f docker/production.yml down -v &> /dev/null
	@echo "Done"
