db-up:
	docker-compose -f docker-compose.postgres.test.yml up -d

migrate:
	migrate -path ./migrations/postgres -database 'postgres://postgres:postgrespw@localhost:5432/postgres?sslmode=disable' up

down:
	migrate -path ./migrations/postgres -database 'postgres://postgres:postgrespw@localhost:5432/postgres?sslmode=disable' down
	docker-compose  down

build: db-up migrate