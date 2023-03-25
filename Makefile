define run-database-prod-migration
	docker run --name my-postgres -e POSTGRES_PASSWORD=postgrespw -e POSTGRES_USER=postgres -e POSTGRES_DATABASE=postgres -p 5432:5432 -d postgres:15
	ping 127.0.0.1 -n 2 > nul
	migrate -path ./migrations/postgres -database 'postgres://postgres:postgrespw@localhost:5432/postgres?sslmode=disable' up 1
endef

test:
	docker compose -f docker-compose.postgres.test.yml up -d
	go test -v ./...
	docker compose -f docker-compose.postgres.test.yml down


run:
	$(call run-database-prod-migration)
	go run main.go;
	docker stop my-postgres

db-up:
	$(call run-database-prod-migration)
