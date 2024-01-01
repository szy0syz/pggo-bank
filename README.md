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

### [Comparing database/sql, GORM, sqlx, and sqlc](https://blog.jetbrains.com/go/2023/04/27/comparing-db-packages/)

```
================================== BENCHMARKING 1 RECORDS ======================================
goos: linux
goarch: amd64
pkg: github.com/rexfordnyrk/go-db-comparison/benchmarks
cpu: Intel(R) Core(TM) i7-8550U CPU @ 1.80GHz
Benchmark/Database/sql_limit:1_-8                   9054            124134 ns/op
Benchmark/Sqlx_limit:1_-8                           8914            138792 ns/op
Benchmark/Sqlc_limit:1_-8                           7954            147056 ns/op
Benchmark/GORM_limit:1_-8                          13388             89251 ns/op
=================================================================================================
================================== BENCHMARKING 10 RECORDS ======================================
Benchmark/Database/sql_limit:10_-8                  7576            157780 ns/op
Benchmark/Sqlx_limit:10_-8                          4384            260402 ns/op
Benchmark/Sqlc_limit:10_-8                          4183            256384 ns/op
Benchmark/GORM_limit:10_-8                          9466            136556 ns/op
=================================================================================================
================================== BENCHMARKING 100 RECORDS ======================================
Benchmark/Database/sql_limit:100_-8                 2521            427603 ns/op
Benchmark/Sqlx_limit:100_-8                         2139            497755 ns/op
Benchmark/Sqlc_limit:100_-8                         2838            456938 ns/op
Benchmark/GORM_limit:100_-8                         1896            563539 ns/op
=================================================================================================
================================== BENCHMARKING 1000 RECORDS ======================================
Benchmark/Database/sql_limit:1000_-8                 516           2201303 ns/op
Benchmark/Sqlx_limit:1000_-8                         445           2786983 ns/op
Benchmark/Sqlc_limit:1000_-8                         535           2313674 ns/op
Benchmark/GORM_limit:1000_-8                         315           4186201 ns/op
=================================================================================================
================================== BENCHMARKING 10000 RECORDS ======================================
Benchmark/Database/sql_limit:10000_-8                 51          21690323 ns/op
Benchmark/Sqlx_limit:10000_-8                         38          28458473 ns/op
Benchmark/Sqlc_limit:10000_-8                         55          21558300 ns/op
Benchmark/GORM_limit:10000_-8                         28          40463924 ns/op
=================================================================================================
================================== BENCHMARKING 15000 RECORDS ======================================
Benchmark/Database/sql_limit:15000_-8                 36          32048808 ns/op
Benchmark/Sqlx_limit:15000_-8                         28          41484578 ns/op
Benchmark/Sqlc_limit:15000_-8                         34          31680017 ns/op
Benchmark/GORM_limit:15000_-8                         20          59348697 ns/op
=================================================================================================
PASS
ok      github.com/rexfordnyrk/go-db-comparison/benchmarks      77.835s
```

#### Conclusion
> While database/sql is the default Golang package, whether you should use it or not depends on your needs as a developer. In this article, you saw the advantages of each package.

> GORM is the best package for you if you need advanced querying, full support functions from underlying database technologies, and a clean codebase—as long as you’re willing to trade off some performance. If you only need basic queries and are comfortable writing your own SQL, then the database/sql or sqlx packages will do just fine.

> Finally, sqlc is best suited for backend developers who work extensively with databases and need to write lots of queries under tight deadlines. You can write your raw SQL queries and generate code without having to worry about types, scanning, or other obstacles to productivity. sqlc also provides great performance gains compared to the rest, especially when dealing with greater volumes of data or record sets.

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

![image](https://github.com/szy0syz/pggo-bank/assets/10555820/5d510f69-b5f6-44c0-88ed-7f83e77defab)

![image](https://github.com/szy0syz/pggo-bank/assets/10555820/5b8d7a6e-96e6-4205-99fd-67483e8d07a0)

![image](https://github.com/szy0syz/pggo-bank/assets/10555820/a4e49613-b46f-4cec-ac01-1601eabed454)

![image](https://github.com/szy0syz/pggo-bank/assets/10555820/d023f09f-2cf2-4a93-96f9-eebd1a6b0afb)

![image](https://github.com/szy0syz/pggo-bank/assets/10555820/122d2ce9-87a4-4f8e-a66d-00e73cc02d24)

### Why mock database?

- Independent tests
- Faster tests
- 100 coverage

How to mock?

Use Fake DB: Memory
implement a fake version of DB: store data in memory

> 为了达到测试覆盖率100%，我们得到所有代码路径全部走一遍，可以为 `黑盒` 准备不同输入，然后遍历这些输入，让其跑满所有路径！🤖

### About params validator

我们在 `json:"currency" binding:"required,currency,oneof=USD EUR CAD"` 做了 `硬编码` 的校验，这里未来有问题，如新增currency类别，不可能重新修改代码。所以这里我们需要custom validator

### Add users table with unique & foreign key constraints

```
migrate create -ext sql -dir db/migration -seq add_users

-- 这两条效果都是一样的
-- 可以用建符合索引的方式来约束出现, user1-USD, user1-USD 的情况
-- CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");
ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");
```

<img width="1314" alt="image" src="https://github.com/szy0syz/pggo-bank/assets/10555820/5e47e728-012a-4256-9d71-9f7b8649444a">

### Why PASETO is better than JWT for token-based authentication?

![image](https://github.com/szy0syz/pggo-bank/assets/10555820/cf26d3d6-85c4-4bdd-84b9-415d8c66743d)

![image](https://github.com/szy0syz/pggo-bank/assets/10555820/81ee1915-8bf6-429f-afc3-26e912986da4)

![image](https://github.com/szy0syz/pggo-bank/assets/10555820/5d009113-1d4b-492c-9059-3640824ddc03)

<img width="767" alt="image" src="https://github.com/szy0syz/pggo-bank/assets/10555820/62d22bf7-aee5-4edc-b99e-62a5e10b57a8">

<img width="881" alt="image" src="https://github.com/szy0syz/pggo-bank/assets/10555820/164cf130-764b-4ac9-8c19-7386fb23b6f1">

- 体会到Go语言面向接口+组合的强大之处 👍🏻 
![image](https://github.com/szy0syz/pggo-bank/assets/10555820/19d90c51-7002-4eb0-9001-8e9b62b6dd94)

### Add new resource

- `make new_migration name=add_sessions`
- update the `db/migration/003_xxx.up.sql` `db/migration/003_xxx.down.sql`
- `make migrateup`
- review database
- update the `db/query/seesion.sql`
- `make sqlc`
- should add `/db/sqlc/session.sql.go` for Golang code
- review the code
- `make mock` regenerate the mock store
- `make test` make sure the all tests is passed

### Database resource

<img width="1444" alt="image" src="https://github.com/szy0syz/pggo-bank/assets/10555820/f0a14208-7c70-475b-99b9-e2630d6c1f86">

- https://dbdocs.io/szy0syz/pggo_bank
- `npm install -g @dbml/cli`
- `make db_schema`
- `make db_docs`

## gRPC

<img width="1358" alt="image" src="https://github.com/szy0syz/pggo-bank/assets/10555820/f3ae1a84-3d94-4277-b227-ad2ed613a525">

- brew install protobuf
- protoc --version
- go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
- go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
- export PATH="$PATH:$(go env GOPATH)/bin"

