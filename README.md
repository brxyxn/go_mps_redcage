# Money Processing System
This project was developed as required by Redcage as a technical test developing and dockerizing a lite version of an API - Rest. Stack used in this project:
- [Golang 1.16](https://go.dev/dl/)
- [PostgreSQL 12.9](https://www.postgresql.org/download/)
- [Docker 20.10.12](https://www.docker.com/get-started)
- [Docker Compose 1.29.2](https://docs.docker.com/compose/)
- [ElephantSQL](https://www.elephantsql.com)
- [Swagger.io](https://swagger.io/)

> Important Note: this project won't accept Pull Requests, however, you can use it as reference if need it to build your own project.

---
# Content
1. [Links](#links)
2. [Getting Started](#getting-started)
	- [Docker](#docker)
		- [Pull repository](#pull-repository)
		- [Versions](#versions)
	- [Source Code](#source-code)
	- [Setting Database Connection](#setting-database-connection)
	- [Initializing The App](#initializing-the-app)
	- [Methods](#methods)
		- [Create a new client](#create-a-new-client)
		- [Create a new account](#create-a-new-account)
		- [Create a new transaction](#create-a-new-transaction)
	- [Database](#database)
3. [What's next?](#whats-next)
---

# Links
<p style="display: flex;flex-direction: row;">
<a href="https://hub.docker.com/repository/docker/brxyxn/go_mps_redcage"><img src="https://upload.wikimedia.org/wikipedia/commons/4/4e/Docker_%28container_engine%29_logo.svg" alt="Docker repository" height="50"></a>

<a href="https://hub.docker.com/repository/docker/brxyxn/go_mps_redcage"><img src="https://upload.wikimedia.org/wikipedia/commons/9/91/Octicons-mark-github.svg" alt="Docker repository" height="50"></a>
</p>


# Getting Started
## Docker
### Pull repository
```bash
docker pull brxyxn/go_mps_redcage:version2
```
### Versions
The `latest` version works with ElephantSQL which is an online platform to connect our PostgreSQL database. However, you can still set up the `dockerized` tag but you must configure the PostgreSQL image locally to be able to run it.

## Source Code
Clone the repository to your local environment.

> HTTPS
```bash
git clone https://github.com/brxyxn/go_mps_redcage.git
```
> SSH
```bash
git clone git@github.com:brxyxn/go_mps_redcage.git
```
> GitHub CLI
```bash
gh repo clone brxyxn/go_mps_redcage
```

## Setting Database Connection
You need to create a local file `.env` with the following information:
```sh
DB_HOST={your_host}
DB_DRIVER={pgx}
DB_USER={your_user}
DB_PASSWORD={your_password}
DB_NAME={backend_db}
DB_PORT={5432}
```

Take into consideration that the information will be needed and you can use the file `.env_template` as a reference.

## Initializing The App
> Move into the code directory
```bash
cd go_mps_redcage
```

> Verify dependencies
```bash
go mod download
go mod tidy
go mod verify
# Output: all modules verified
```

> Run
```bash
go run .
# or
go run main.go
# Output:
# go-mps-api 2022/03/29 23:53:57 [INFO] Running server on port :3000
```
> After running the app you will see the output showing that the app is initialized and running.


If the tables do not exist in the database the app will automatically create them. Refer below to see the schema. By default every table will be dropped for testing purposes, but you can remove those SQL queries from the file `docker_postgres_init.sql`.

```plsql
-- Remove these lines if you need to keep records up
DROP TABLE IF EXISTS public.transactions;
DROP TABLE IF EXISTS public.accounts;
DROP TABLE IF EXISTS public.clients;
```

## Methods

| HTTP Method | URI Pattern                                                            |
| -----------:| ---------------------------------------------------------------------- |
| POST        | /api/v1/clients                                                        |
| GET         | /api/v1/clients/**{client_id}**                                        |
| POST        | /api/v1/clients/**{client_id}**/accounts                               |
| GET         | /api/v1/clients/**{client_id}**/accounts                               |
| GET         | /api/v1/clients/**{client_id}**/accounts/**{account_id}**              |
| POST        | /api/v1/clients/**{client_id}**/accounts/**{account_id}**/transactions |
| GET         | /api/v1/clients/**{client_id}**/accounts/**{account_id}**/transactions |

### Create a new client
> /api/v1/clients

Accepted fields:
- firstName as **string**.
- lastName as **string**.
- username as **string**.

``` json
{
	"firstName": "John",
	"lastName": "Doe",
	"username": "john.doe"
}
```

### Create a new account
> /api/v1/clients/*1*/accounts

Accepted fields:
- balance as **numeric(13,2)**.  
- currency codes as **string**:
  - USD
  - MXN
  - COP
- accountType as **string**:
  - Savings
  - Checking
  - Credit Card

For more details about currency symbols and codes check [this link.](https://www.xe.com/symbols.php)

```json
{
	"balance": 80000.00,
	"currency": "USD",
	"accountType": "Savings"
}
```
```json
{
	"balance": 25000.00,
	"currency": "MXN",
	"accountType": "Checking"
}
```
```json
{
	"balance": 65000.00,
	"currency": "COP",
	"accountType": "Credit Card"
}

```
### Create a new transaction
> /api/v1/clients/*5*/accounts/*8*/transactions

Accepted fields:
- amount as **numeric(13,2)**.  
- transactionType as **integer**:
  - 1 -> Deposit
  - 2 -> Withdraw
  - 3 -> Transfer
- description as **string**.  
- receiverAccountId as **integer**.  
- senderAccountId as **integer**.


>If the value for senderAccountId is not set, by default will be 0 which means the transaction is either a deposit or a withdrawal.

```json
{
	"amount": 100.00,
	"transactionType": 1,
	"description": "This is a deposit",
	"receiverAccountId": 1,
	"senderAccountId": 0
}
```
```json
{
	"amount": 120.00,
	"transactionType": 2,
	"description": "This is a withdraw",
	"receiverAccountId": 2,
	"senderAccountId": 0
}
```
```json
{
	"amount": 100.00,
	"transactionType": 3,
	"description": "This is a transfer or EFT",
	"receiverAccountId": 2,
	"senderAccountId": 1
}
```

## Database
> This is the final database schema, remember to update the public schema when deployed to production.  
> You can either run the script manually into the database or let the API do it for you.

```sql
CREATE TABLE IF NOT EXISTS public.client (
	id integer NOT NULL GENERATED ALWAYS AS IDENTITY,
	firstname varchar NOT NULL,
	lastname varchar NOT NULL,
	username varchar NOT NULL UNIQUE,
	active boolean NOT NULL DEFAULT TRUE,
	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	
	CONSTRAINT client_pk PRIMARY KEY (id),
	CONSTRAINT client_username_key UNIQUE (username)
);

CREATE TABLE IF NOT EXISTS public.account (
	id integer NOT NULL GENERATED ALWAYS AS IDENTITY,
	balance numeric(13,2) NOT NULL,
	currency varchar NOT NULL,
	account_type varchar NOT NULL,
	active boolean NOT NULL DEFAULT True,
	client_id integer NOT NULL,
	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	
	CONSTRAINT account_pk PRIMARY KEY (id),
	CONSTRAINT account_fk FOREIGN KEY (client_id) REFERENCES public.client(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS public."transaction" (
	id integer NOT NULL GENERATED ALWAYS AS IDENTITY,
	amount numeric(13,2) NOT NULL,
	transaction_type smallint NOT NULL,
	description varchar NOT NULL UNIQUE,
	receiver_account_id integer NOT NULL,
	sender_account_id integer,
	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,

	CONSTRAINT transaction_pk PRIMARY KEY (id),
	CONSTRAINT transaction_fk FOREIGN KEY (receiver_account_id) REFERENCES public.account(id) ON DELETE RESTRICT
);

```
> A SQL file is included with dummy data to get easy to start working with the API, you could start creating transactions which is the main target of this project.

---

# What's next?
For the following development and releases there are some ideas of what could be implemented, eg: creating a profile for clients where additional "real-life" information is collected.
```golang
type Profile struct {
	Profile           uint64      `json:"profileId"`
	DateOfBirth       string      `json:"dateOfBirth"`
	ProfilePicture    string      `json:"profilePicture"` // base64
	Contact           ContactInfo `json:"contact"`
	PhysicalAddress   Address     `json:"physicalAddress"`
	MailingAddress    Address     `json:"mailingAddress"`
	AdditionalDetails string      `json:"additionalDetails"`
	ClientId          uint64      `json:"clientId"`
}

type ContactInfo struct {
	PrimaryEmail         string      `json:"primaryEmail"`
	SecondaryEmail       string      `json:"secondaryEmail"`
	PrimaryPhoneNumber   PhoneNumber `json:"primaryPhoneNumber"`
	SecondaryPhoneNumber PhoneNumber `json:"secondaryPhoneNumber"`
}

type PhoneNumber struct {
	AreaCode          string `json:"areaCode"`
	Number            string `json:"number"`
	ContactByText     bool   `json:"contactByText"`
	ContactByVoice    bool   `json:"contactByVoice"`
	ContactByWhatsApp bool   `json:"contactByWhatsApp"`
}

type Address struct {
	AddressLineOne string `json:"addressLineOne"`
	AddressLineTwo string `json:"addressLineTwo"`
	City           string `json:"city"`
	State          string `json:"state"`
	Country        string `json:"country"`
	ClientId       uint64 `json:"clientId"`
}
```

> Additionally there are some improvements based on the "fictitious" client's requirement, however, those might be covered at another time.
