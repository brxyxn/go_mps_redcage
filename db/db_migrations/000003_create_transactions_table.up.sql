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