include .envrc


.PHONY: run/api
run/api:
	go run ./cmd/bot -db-dsn=${TELEGRAM_BOT}?sslmode=disable


.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]


.PHONY: start-container
start-container:
	docker run --name telegram-bot-v03 -p 8080:8080 --env-file .env telegram-bot:0.1

.PHONY: build-image
build-image:
	docker build -t telegram-bot-v03 .

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${TELEGRAM_BOT}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${TELEGRAM_BOT} up

current_time = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
linker_flags = '-s -X main.buildTime=${current_time}'

.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags=${linker_flags} -o=./bin/bot ./cmd/bot
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/bot ./cmd/bot
