-- DB name is jwt_users
-- create users table
-- and
-- create jwt table fk to user table
-- and
-- secrets table

BEGIN;
CREATE DATABASE jwt_users;
COMMIT;

BEGIN;
CREATE TABLE IF NOT EXISTS jwt_users.users (
    id serial PRIMARY KEY,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    location TEXT,
    ip_addr inet,
    created_at timestamp not null default now()
);

CREATE TABLE IF NOT EXISTS jwt_users.tokens (
    user_id integer references users(id) on delete cascade,
    jwt_token TEXT NOT NULL,
    created_at timestamp,
    expires_at timestamp
);

Create table IF NOT EXISTS jwt_users.secrets (
    project_name TEXT primary Key,
    secret_key TEXT NOT NULL,
    created_at timestamp not null default now(),
    updated_at timestamp
);
COMMIT;

-- lets get a different user than postgres to own everything. something more specific to the proj
begin;
CREATE ROLE token_master with login;
commit;

begin;
alter database jwt_users owner to token_master;
alter table jwt_users.users owner to token_master;
alter table jwt_users.tokens owner to token_master;
alter table jwt_users.secrets owner to token_master;
commit;

-- changed the users and tokens table, so had to drop them
-- had to recreate, but then i also have to reapply ownerships
-- so creatign a schema, adding tables to that schema and letting token_master own schema
begin;
create schema if not exists jwt_auth;
alter schema jwt_auth owner to token_master;
ALTER ROLE token_master SET search_path TO jwt_auth;
alter table jwt_users.users set schema jwt_auth;
alter table jwt_users.tokens set schema jwt_auth;
alter table jwt_users.secrets set schema jwt_auth;

ALTER ROLE token_master SET search_path TO jwt_auth;
commit;
begin;
alter table jwt_users.users add constraint unique_username unique (username);
commit;