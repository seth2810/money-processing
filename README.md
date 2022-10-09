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

The following methods are implemented:

- Create Client (POST /clients)
- Create Account for Client (POST /accounts)
- Get Client (GET /clients/:id)
- Get Account (GET /accounts/:id)
- Get Transactions - return list of transactions for account (GET /accounts/:id/transactions)
- Create Transaction - create transaction of needed type (POST /transactions)

API - JSON via REST
Programming language - Golang
Database - PostgreSQL
