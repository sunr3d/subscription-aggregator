env:
	if [ ! -f .env ]; then cp .env.example .env; fi

up: env
	docker compose up -d --build

down:
	docker compose down

restart: down up

clean:
	docker compose down -v

logs:
	docker compose logs -f app

test:
	go test -v ./...