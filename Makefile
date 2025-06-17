.PHONY: all up down build migrate-up migrate-down migrate-create seed logs ps psql redis clean

include .env
export $(shell sed 's/=.*//' .env)

MIGRATE_DATABASE_URL = postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@db:5432/${POSTGRES_DB}?sslmode=disable

all: build up migrate-up seed

up:
	docker-compose up -d

down:
	docker-compose down -v

build:
	docker-compose build --no-cache
migrate-up:
	docker-compose exec app migrate -path /app/insider-project/migrations -database "postgres://postgres:secret@db:5432/mydb?sslmode=disable" up
migrate-down:
	docker-compose exec app sh -c "cd /app/insider-project && migrate -path ./migrations -database ${MIGRATE_DATABASE_URL} down 1"

migrate-create: build
	docker-compose exec app sh -c "cd /app/insider-project && migrate create -ext sql -dir ./migrations ${name}"

seed:
	docker-compose exec app go run ./cmd/seed/seed_runner.go

logs:
	docker-compose logs -f app

ps:
	docker-compose ps

psql:
	docker-compose exec db /usr/bin/psql -U ${POSTGRES_USER} -d ${POSTGRES_DB}

redis:
	docker-compose exec redis redis-cli

clean: down
	docker volume prune -f
	docker image prune -a -f
