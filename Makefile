include .env
export

export PROJECT_ROOT= $(shell pwd)

.PHONY: db-up db-down db-clean db-logs migrate-create migrate-command migrate-up migrate-down

# Database commands
db-up:
	@docker compose up -d

db-down:
	@docker compose down

db-clean:
	@read -p "Remove all data from the database? (y/N): " confirm; \
	if [ "$$confirm" = "y" ]; then \
		docker compose down -v && \
		rm -rf out/pgdata && \
		echo "Database data removed"; \
	fi

db-logs:
	@docker compose logs -f

# Create a new migration
migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Error: sequence number is required"; \
		exit 1; \
	fi

	@docker compose run --rm migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

# Run a custom migration command
migrate-command:
	@if [ -z "$(command)" ]; then \
		echo "Error: command is required"; \
		exit 1; \
	fi

	@docker compose run --rm migrate \
		-path /migrations \
		-database postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@db:5432/$(POSTGRES_DB)?sslmode=disable \
		"$(command)"

# Run the up migration
migrate-up:
	@make migrate-command command=up

# Run the down migration
migrate-down:
	@make migrate-command command=down