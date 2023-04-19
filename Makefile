INTEGRATION_TEST_PATH?=./pkg/repository/

build:
	docker compose -f docker-compose.yml build

#make run file=
run:
	docker compose -f $(file) up -d

down:
	docker compose -f $(file) down


lint:
	golangci-lint run ./...

#make add name=
add:
	migrate create -ext sql -dir ./internal/migrations/postgres -seq $(name)



test.repository:
	docker compose -f docker-compose.test.repository.postgres.yml up -d
	go test -tags=integration $(INTEGRATION_TEST_PATH) -count=1 -v
	docker compose -f docker-compose.test.repository.postgres.yml down




