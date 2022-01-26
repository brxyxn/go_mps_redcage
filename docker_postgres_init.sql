DROP TABLE IF EXISTS public."transaction";
DROP TABLE IF EXISTS public.account;
DROP TABLE IF EXISTS public.client;

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

-- Example data
-- Clients
insert into public.client(firstname, lastname, username) values('John', 'Doe', 'supecooluser1');
insert into public.client(firstname, lastname, username) values('Hugo', 'First', 'supecooluser2');
insert into public.client(firstname, lastname, username) values('Ray', 'Sin', 'supecooluser3');
insert into public.client(firstname, lastname, username) values('Dee', 'End', 'supecooluser4');
insert into public.client(firstname, lastname, username) values('Karen', 'Fays', 'supecooluser5');
insert into public.client(firstname, lastname, username) values('Anne', 'Dote', 'supecooluser6');
-- Accounts
insert into public.account(balance, currency, account_type, client_id) values (0.00, 'USD', 'Savings', 1);
insert into public.account(balance, currency, account_type, client_id) values (0.00, 'MXN', 'Savings', 1);
insert into public.account(balance, currency, account_type, client_id) values (0.00, 'COP', 'Savings', 1);
insert into public.account(balance, currency, account_type, client_id) values (0.00, 'USD', 'Checking', 2);
insert into public.account(balance, currency, account_type, client_id) values (0.00, 'MXN', 'Checking', 2);
insert into public.account(balance, currency, account_type, client_id) values (0.00, 'COP', 'Checking', 2);
insert into public.account(balance, currency, account_type, client_id) values (0.00, 'USD', 'Savings', 3);
insert into public.account(balance, currency, account_type, client_id) values (0.00, 'MXN', 'Savings', 3);
insert into public.account(balance, currency, account_type, client_id) values (0.00, 'COP', 'Savings', 3);
insert into public.account(balance, currency, account_type, client_id) values (0.00, 'USD', 'Checking', 4);
insert into public.account(balance, currency, account_type, client_id) values (0.00, 'MXN', 'Checking', 4);
insert into public.account(balance, currency, account_type, client_id) values (0.00, 'COP', 'Checking', 4);
insert into public.account(balance, currency, account_type, client_id) values (0.00, 'USD', 'Savings', 5);
insert into public.account(balance, currency, account_type, client_id) values (0.00, 'MXN', 'Savings', 5);
insert into public.account(balance, currency, account_type, client_id) values (0.00, 'COP', 'Savings', 5);
insert into public.account(balance, currency, account_type, client_id) values (0.00, 'USD', 'Credit Card', 6);
insert into public.account(balance, currency, account_type, client_id) values (0.00, 'MXN', 'Credit Card', 6);
insert into public.account(balance, currency, account_type, client_id) values (0.00, 'COP', 'Credit Card', 6);