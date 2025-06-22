.PHONY: up down logs clean trash setup

up:
	docker compose -f docker-compose.yml up -d

down:
	docker-compose -f docker-compose.yml down

logs:
	docker-compose -f docker-compose.yml logs -f

clean:
	docker compose -f docker-compose.yml down -v

trash:
	docker compose -f docker-compose.yml down -v --rmi all

setup:
	bash scripts/setup.sh

seed:
	curl -X GET http://localhost:8080/seed
