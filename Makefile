.PHONY: up down clean install-swag swagger run-e2e-tests
PKGS := $(shell go list ./... | grep -vE 'mocks|cmd/pr-review')

up:
	docker compose up --build -d

down:
	docker compose down

clean:
	docker compose down -v

install-swag:
	go install github.com/swaggo/swag/cmd/swag@latest

swagger:
	swag init \
  -d ./cmd/pr-review,./internal/http/server,./internal/http/handlers \
  --parseDependency \
  --parseInternal \
  -o ./docs

run-e2e-tests:
	docker compose -f docker-compose.tests.yml up -d --build db
	until [ "`docker compose -f docker-compose.tests.yml ps -q db | xargs docker inspect -f '{{ .State.Health.Status }}'`" = "healthy" ]; do \
		echo "Waiting for DB..."; sleep 1; \
	done
	- go test -v ./e2e/...
	docker compose -f docker-compose.tests.yml down -v

coverage:
	docker compose -f docker-compose.tests.yml up -d --build db
	until [ "`docker compose -f docker-compose.tests.yml ps -q db | xargs docker inspect -f '{{ .State.Health.Status }}'`" = "healthy" ]; do \
		echo "Waiting for DB..."; sleep 1; \
	done
	go test $(PKGS) -coverprofile=coverage.out -coverpkg=$(shell echo $(PKGS) | tr ' ' ',') -timeout 30s
	go tool cover -html=coverage.out -o coverage.html
	docker compose -f docker-compose.tests.yml down -v
