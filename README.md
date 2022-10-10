# Money Processing Service

The test is to implement a light version of a bank processing service.

Just like a typical bank, it operates following entities:

**Clients** of the bank who keep money there.
Each **Client** can have multiple **Accounts**. Each **Account** balance can be denominated in one currency.
For example Client1 can have three accounts: one USD account and two COP accounts.
The following currencies should be supported: USD, COP, MXN.
Each Client can have multiple accounts with the same currency.

**Transaction** is an action to update an **Account** or **Accounts** balance.
A transaction should belong to one of following types:

- Deposit money
- Withdraw money
- Money transfer between Accounts (Currency conversion is not needed. Only transfer between accounts with the same currency is allowed).

## Quick start

To start server, you basically just have to start database and server docker containers and you are ready to go:

```
$ docker run -d \
    --name money-processing-db \
    -e POSTGRES_USER=money-processing \
    -e POSTGRES_PASSWORD=money-processing \
    -e POSTGRES_DB=money-processing \
    -p 5432:5432 \
    postgres:14.5-alpine

$ docker run -d \
    --platform linux/amd64 \
    --name money-processing-server \
    --link money-processing-db:db \
    -e DB_HOST=db \
    -e DB_PORT=5432 \
    -e DB_USER=money-processing \
    -e DB_PASSWORD=money-processing \
    -e DB_NAME=money-processing \
    -e HOST=0.0.0.0 \
    -e PORT=8080 \
    -p 8080:8080 \
    seth2810/money_processing:develop
```

or via `docker compose` inside project root directory:

```
$ docker compose up -d
# or via make
$ make up
```

Then open separate terminal and run API methods.

## API

The following methods are implemented:

### Create clients

```
$ curl -X POST http://localhost:8080/clients \
   -H 'Content-Type: application/json' \
   -d '{"email":"john@example.com"}'

{"id":1,"email":"john@example.com"}

$ curl -X POST http://localhost:8080/clients \
   -H 'Content-Type: application/json' \
   -d '{"email":"sam@example.com"}'

{"id":2,"email":"sam@example.com"}
```

### Create accounts

```
$ curl -X POST http://localhost:8080/accounts \
   -H 'Content-Type: application/json' \
   -d '{"client_id": 1, "currency":"USD"}'

{"id":1,"client_id":1,"balance":"0","currency":"USD"}

$ curl -X POST http://localhost:8080/accounts \
   -H 'Content-Type: application/json' \
   -d '{"client_id": 2, "currency":"USD"}'

{"id":2,"client_id":2,"balance":"0","currency":"USD"}
```

### Get client

```
$ curl http://localhost:8080/clients/1

{"id":1,"email":"john@example.com"}
```

### Get account

```
$ curl http://localhost:8080/accounts/2

{"id":2,"client_id":2,"balance":"0","currency":"USD"}
```

### Create Transactions

```
$ curl -X POST http://localhost:8080/transactions \
   -H 'Content-Type: application/json' \
   -d '{"type": "deposit", "to_account_id": 1, "amount": "65.4321"}'

{"id":1,"type":"deposit","amount":"65.4321","from_account_id":null,"to_account_id":1,"created_at":"2022-10-10T12:18:21.733171Z"}

$ curl -X POST http://localhost:8080/transactions \
   -H 'Content-Type: application/json' \
   -d '{"type": "transfer", "from_account_id": 1, "to_account_id": 2, "amount": "23.45"}'

{"id":2,"type":"transfer","amount":"23.45","from_account_id":1,"to_account_id":2,"created_at":"2022-10-10T12:19:58.041629Z"}

$ curl -X POST http://localhost:8080/transactions \
   -H 'Content-Type: application/json' \
   -d '{"type": "withdraw", "from_account_id": 1, "amount": "23.45"}'

{"id":3,"type":"withdraw","amount":"23.45","from_account_id":1,"to_account_id":null,"created_at":"2022-10-10T12:21:43.902788Z"}
```

### Get Transactions - return list of transactions for account

```
$ curl http://localhost:8080/accounts/1/transactions

[{"id":1,"type":"deposit","amount":"65.4321","from_account_id":null,"to_account_id":1,"created_at":"2022-10-10T12:18:21.733171Z"},{"id":2,"type":"transfer","amount":"23.45","from_account_id":1,"to_account_id":2,"created_at":"2022-10-10T12:19:58.041629Z"},{"id":3,"type":"withdraw","amount":"23.45","from_account_id":1,"to_account_id":null,"created_at":"2022-10-10T12:21:43.902788Z"}]

$ curl http://localhost:8080/accounts/2/transactions

[{"id":2,"type":"transfer","amount":"23.45","from_account_id":1,"to_account_id":2,"created_at":"2022-10-10T12:19:58.041629Z"}]
```

## Technical specifications

- Database - PostgreSQL
- Programming language - Golang
- API - JSON via REST
