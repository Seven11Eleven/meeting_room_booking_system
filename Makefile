.PHONY: test

test:
	@echo "Starting test database and migration..."
	@docker compose -f docker-compose.test.yaml up -d db_test migrate
	@sleep 5 
	@echo "Running tests..."
	@go test ./internal/services -v
	@echo "Stopping test database..."
	@docker compose -f docker-compose.test.yaml down
