up:
	docker compose up --build -d

down:
	docker compose down

clean:
	docker compose down -v

install-swag:
	go install github.com/swaggo/swag/cmd/swag@latest

swagger:
	swag init -g ./internal/http/server/router.go --parseDependency --parseInternal --generatedTime -o ./docs
