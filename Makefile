.PHONY: up down migrate

up:
	docker-compose up -d

down:
	docker-compose down

migrate:
	docker exec -it notes-api-postgres psql -U postgres -d notes_api_db -c "\\dt"