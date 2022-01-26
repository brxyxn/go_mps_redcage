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