.PHONY: run
run:
	REDIS_ADDR=localhost:6379 GATEWAY_ADDR=http://localhost:8888 POSTGRES_URL="postgres://user:password@localhost:5432/db?sslmode=disable" air -c .air.toml

.PHONY: up
up:
	docker compose up

.PHONY: update
update:
	@go mod tidy && go get -u ./...
