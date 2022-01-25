# Money Processing System
This project was developed as required by Redcage as a technical test developing and dockerizing a lite version of an API - Rest. Stack used in this project:
- [Golang 1.16](https://go.dev/dl/)
- [PostgreSQL 12.9](https://www.postgresql.org/download/)
- [Docker 20.10.12](https://www.docker.com/get-started)

## Methods
|HTTP Method|URI Pattern |
|---:|---|
|POST|/api/v1/client/new|
|GET|/api/v1/client/**{client_id}**|
|POST|/api/v1/client/**{client_id}**/accounts/new|
|GET|/api/v1/client/**{client_id}**/accounts/**{account_id}**|
|POST|/api/v1/client/**{client_id}**/accounts/**{account_id}**/transactions/new|
|GET|/api/v1/client/**{client_id}**/accounts/**{account_id}**/transactions|

# Getting Started
## Source Code
Clone the repository to your local environment.

### HTTPS
```bash
git clone https://github.com/brxyxn/go_mps_redcage.git
```
### SSH
```bash
git clone git@github.com:brxyxn/go_mps_redcage.git
```
### GitHub CLI
```bash
gh repo clone brxyxn/go_mps_redcage
```

## Running The App
```bash
go run .
```
```bash
go run main.go
```
After running the app you will see the output showing that the app is initialized and running.
```bash
# Output
2022/01/25 12:57:22 API initialized!
2022/01/25 12:57:22 Running at localhost:5000
```
If the tables not exist in the database the app will automatically create them. Refer below to see the schema.



## Database
This is the final database schema, remember to update the public schema when deployed to production.

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