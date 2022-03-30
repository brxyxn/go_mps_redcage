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