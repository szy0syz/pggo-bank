DB_URL=postgresql://root:postgres@localhost:5430/pggo_bank?sslmode=disable
DB_CONTAINER=bank-postgres

network:
	docker network create bank-network

postgres:
	docker run --name "$(DB_CONTAINER)" --network bank-network -p 5430:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=postgres -d postgres:14-alpine

createdb:
	docker exec -it "$(DB_CONTAINER)" createdb --username=root --owner=root pggo_bank

dropdb:
	docker exec -it postgres dropdb pggo_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/szy0syz/pggo-bank/db/sqlc Store

.PHONY: network postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 new_migration db_docs db_schema sqlc test server mock proto evans redis