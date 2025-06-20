-- Create a database names jwt_users;
-- CREATE DATABASE jwt_users;
-- from the cmdline do:
-- sudo psql -U postgres -d jwt_users -f <this_files_name>

-- the login after setup is easy for mac, but gotta jump through some hoops for linux. 
-- you can either:
-- -- mess around with some pg_hba.conf files
-- -- or create a system user with names token_master
-- -- -- sudo adduser token_master
-- -- -- then log in as: sudo -u token_master psql -d jwt_users
-- better than the alternative of:
-- sudo -u postgres psql
-- \c jwt_users;
-- set role token_master;
-- set search_path to jwt_auth;

-- IF YOU NEED TO DROP EVERYTHING:
-- Drop database jwt_users;
-- drop role token_master;

-- Tables are in a schema. Schema is in a database. Roles are system level

BEGIN;
CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    location TEXT,
    ip_addr inet,
    created_at timestamp not null default now()
);

CREATE TABLE IF NOT EXISTS tokens (
    user_id integer references users(id) on delete cascade,
    jwt_token TEXT NOT NULL,
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


begin;
alter database jwt_users owner to token_master;
alter table users owner to token_master;
alter table tokens owner to token_master;
alter table secrets owner to token_master;
commit;

-- changed the users and tokens table, so had to drop them
-- had to recreate, but then i also have to reapply ownerships
-- so creatign a schema, adding tables to that schema and letting token_master own schema
begin;
create schema if not exists jwt_auth;
alter schema jwt_auth owner to token_master;
ALTER ROLE token_master SET search_path TO jwt_auth;
alter table users set schema jwt_auth;
alter table tokens set schema jwt_auth;
alter table secrets set schema jwt_auth;
commit;

begin;
alter table jwt_auth.users add constraint unique_username unique (username);
commit;