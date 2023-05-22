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

ACID

- Atomicity (A)
  - Either all operations complete successfully or then transaction fails and the db is unchanged.
- Consistency (C)
  - The db state must be valid after the transaction. All constraints must be satisfied.
- Isolation (I)
  - Concurrent trans must not affect each other.
- Durability (D)
  - Data written by a successful transaction must be recorded in persistent storage

`FOR Update`

```
-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1
FOR UPDATE;
```

- GetAccount 是一个普通的 SELECT 查询，它只会读取匹配的数据，并不会对数据进行锁定。
- GetAccountForUpdate 同样是 SELECT 查询，但是它带有 FOR UPDATE 子句，它会对查询到的数据行应用一个排他锁（Exclusive Lock）。这意味着，直到当前事务结束，其他事务都无法修改这些被锁定的数据行，以保证数据的一致性和防止并发事务产生的数据竞态问题。这在需要进行复杂更新或者是保证在一系列操作中数据不会被其他事务更改的情况下非常有用。

需要注意的是，FOR UPDATE 锁只在事务内有效，因此这个查询通常需要在 BEGIN 和 COMMIT / ROLLBACK 语句之间执行。否则，由于大多数数据库系统在查询结束后自动提交事务，FOR UPDATE 锁会立即被释放。

所以，总的来说，这两句SQL查询的主要区别在于后者会锁定选定的数据行，防止其他事务对这些行进行修改，直到当前事务结束。

DB transaction lock & How to handle deadlock

- 首先发现外键导致了死锁，取消外键可解决，但这无法满足数据一致性要求
- 然后再goroutine里打印日志逐步排查
- 最后使用`FOR NO KEY UPDATE;`解决
- 最终使用一步更新，不Select锁定优化

<img width="877" alt="image" src="https://github.com/szy0syz/pggo-bank/assets/10555820/f8fa655d-8692-4197-8e21-0a930e380aa1">

### How to avoid deadlock in DB transaction?

> Queries order matters!

```sql
BEGIN;

UPDATE accounts SET balance = balance - 10 WHERE "id" = 1 RETURNING *;
UPDATE accounts SET balance = balance + 10 WHERE "id" = 2 RETURNING *;

ROLLBACK;



-- Tx2: transfer $10 from account2 to account1
BEGIN;

UPDATE accounts SET balance = balance - 10 WHERE "id" = 2 RETURNING *;
UPDATE accounts SET balance = balance + 10 WHERE "id" = 1 RETURNING *;

ROLLBACK;

-- docker exec -it bank-postgres psql -U root -d pggo_bank
```

```
pggo_bank=*# UPDATE accounts SET balance = balance + 10 WHERE "id" = 1 RETURNING *;
ERROR:  deadlock detected
DETAIL:  Process 60 waits for ShareLock on transaction 841; blocked by process 50.
Process 50 waits for ShareLock on transaction 842; blocked by process 60.
HINT:  See server log for query details.
CONTEXT:  while updating tuple (0,74) in relation "accounts"
```

### Understand isolation levels & read phenomena


![image](https://github.com/szy0syz/pggo-bank/assets/10555820/d023f09f-2cf2-4a93-96f9-eebd1a6b0afb)

![image](https://github.com/szy0syz/pggo-bank/assets/10555820/122d2ce9-87a4-4f8e-a66d-00e73cc02d24)
