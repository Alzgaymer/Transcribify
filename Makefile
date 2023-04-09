up:
	docker compose -f docker-compose.yml build
	docker compose -f docker-compose.yml up -d

down:
	docker compose -f docker-compose.yml down


lint:
	golangci-lint run ./...

