build:
	docker compose -f docker-compose.yml build

run:
	docker compose -f docker-compose.yml up -d

down:
	docker compose -f docker-compose.yml down


lint:
	golangci-lint run ./...

add:
	migrate create -ext sql -dir ./internal/migrations/postgres -seq $(name)

