include .envrc

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	@echo 'Running application…'
	@go run ./cmd/api \
		-port=4000 \
		-env=development \
		-limiter-burst=5 \
		-limiter-rps=2 \
		-limiter-enabled=true \
		-db-dsn=${ANIMEVERSE_DB_DSN} \
		-cors-trusted-origins="http://localhost:8081 http://localhost:9000 http://localhost:9001"

## db/psql: connect to the database using psql (terminal)
.PHONY: db/psql
db/psql:
	psql ${ANIMEVERSE_DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${ANIMEVERSE_DB_DSN} up

## db/migrations/down: apply all down database migrations
.PHONY: db/migrations/down
db/migrations/down:
	@echo 'Running down migrations...'
	migrate -path ./migrations -database ${ANIMEVERSE_DB_DSN} down

## db/migrations/force version=$1: force database migration version
.PHONY: db/migrations/force
db/migrations/force:
	@echo 'Forcing migration version to ${version}...'
	migrate -path ./migrations -database ${ANIMEVERSE_DB_DSN} force ${version}


