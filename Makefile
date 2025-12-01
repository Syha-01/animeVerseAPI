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

## test: run all tests
.PHONY: test
test:
	@echo 'Running tests...'
	@go test ./...

## test/verbose: run all tests with verbose output
.PHONY: test/verbose
test/verbose:
	@echo 'Running tests (verbose)...'
	@go test -v ./...

## test/coverage: run all tests with coverage
.PHONY: test/coverage
test/coverage:
	@echo 'Running tests with coverage...'
	@go test -cover ./...

## test/coverage-report: generate HTML coverage report
.PHONY: test/coverage-report
test/coverage-report:
	@echo 'Generating coverage report...'
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo 'Coverage report generated: coverage.html'


