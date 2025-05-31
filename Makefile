DB_URL=postgresql://root:cocoa@localhost:5432/alabast?sslmode=disable

postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=cocoa -d postgres:16.3

createdb:
	docker exec -it postgres createdb --username=root --owner=root alabast

dropdb:
	docker exec -it postgres dropdb alabast

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

sqlc:
	sqlc generate

mock:
	mockgen -package mockdb -destination db/mock/store.go alabast/db/sqlc Store

test:
	go test -v -cover ./...

server:
	go run main.go

