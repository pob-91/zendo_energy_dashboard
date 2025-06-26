.PHONY: start up down logs clean trash

start:
	if [ -z "$$(docker images -q zendo-data-fetcher:latest)" ]; then cd data_fetcher && bash build.sh; fi
	if [ -z "$$(docker images -q zendo-api:latest)" ]; then cd api && bash build.sh; fi
	if [ -z "$$(docker images -q zendo-web-app:latest)" ]; then cd web_client && bash build.sh; fi
	if [ -z "$$(docker images -q zendo-data-processor:latest)" ]; then cd data_processor && bash build.sh; fi
	if [ -z "$$(docker images -q zendo-cron:latest)" ]; then cd cron && bash build.sh; fi
	docker compose -f docker-compose.yml up -d main-db data-fetcher
	sleep 5 # this is not great but docker compose waiting is a pain
	bash scripts/setup.sh
	curl -X GET http://localhost:8080/seed && cd data_processor && uv run process_historical.py && cd ..
	docker compose -f docker-compose.yml up -d

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
	@if [ -n "$$(docker images -q golang:alpine)" ]; then docker rmi golang:alpine; fi
	@if [ -n "$$(docker images -q node:jod-alpine)" ]; then docker rmi node:jod-alpine; fi
