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
