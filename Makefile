postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USERNAME=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank2

dropdb:
	docker exec -it postgres dropdb simple_bank2

sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb sqlc 
