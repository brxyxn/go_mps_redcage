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