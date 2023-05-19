# pggo-bank
pg-go-k8s-grpc

## Notes

```
migrate create -ext sql -dir db/migration -seq init_schema

make sqlc
```

Database/SQL

- Very fast & straightforward
- Manual mapping SQL fields to variables
- Easy to make mistakes, not caught until runtime

GORM

- CURD functions already implemented, very short production code
- Must learn to write queries using gorm's function
- Run slowly on high load

SQLX

- Quite fast & easy to use
- Fields mapping via query text & struct tags
- Failure won't occur until runtime

SQLC

- Very fast & easy to use
- Automatic code generation
- Catch SQL query errors before generating codes