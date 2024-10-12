db:
	docker run --name postgres -p 8015:5432 -e POSTGRES_PASSWORD=dev -d --restart=always --network app postgres:15.6
d_build:
	docker build -t astroservice:local .

d_run:
	docker rm -f astroservice && docker run --name astroservice --network app -p 8088:8088 -v $(pwd)/.env-docker:/app/.env:ro -d astroservice:local

d_up:
	docker compose up -d

test:
	docker rm -f postgrestest
	docker run --name postgrestest -p 8091:5432 -e POSTGRES_PASSWORD=dev -e POSTGRES_DB=postgres -d postgres:15.6
	sleep 3
	goose -dir ./migrations postgres "postgres://postgres:dev@localhost:8091/postgres?sslmode=disable" up
	go test -v -race ./...

lint:
	golangci-lint  run --enable-all