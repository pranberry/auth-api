-- create users table
-- and
-- create jwt table fk to user table


CREATE TABLE users (
    username TEXT NOT NULL,
    id serial PRIMARY KEY,
    password TEXT NOT NULL,
    location TEXT,
    ip_addr inet,
    created_at timestamp not null default now()
);

CREATE TABLE tokens (
    jwt_token TEXT NOT NULL,
    user_id integer references users(id) on delete cascade,
    created_at timestamp,
    expires_at timestamp
);

Create table secrets (
    project_name TEXT primary Key,
    secret_key TEXT NOT NULL
    created_at timestamp not null default now()
    updated_at timestamp
);