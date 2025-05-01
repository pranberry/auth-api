-- DB name is jwt_users
-- create users table
-- and
-- create jwt table fk to user table
-- and
-- secrets table

BEGIN;
CREATE TABLE IF NOT EXISTS users (
    username TEXT NOT NULL,
    id serial PRIMARY KEY,
    password TEXT NOT NULL,
    location TEXT,
    ip_addr inet,
    created_at timestamp not null default now()
);

CREATE TABLE IF NOT EXISTS tokens (
    jwt_token TEXT NOT NULL,
    user_id integer references users(id) on delete cascade,
    created_at timestamp,
    expires_at timestamp
);

Create table IF NOT EXISTS secrets (
    project_name TEXT primary Key,
    secret_key TEXT NOT NULL,
    created_at timestamp not null default now(),
    updated_at timestamp
);
COMMIT;

-- lets get a different user than postgres to own everything. something more specific to the proj
begin;
CREATE ROLE token_master with login;
alter database jwt_users owner to token_master;
alter table users owner to token_master;
alter table tokens owner to token_master;
alter table secrets owner to token_master;
commit;