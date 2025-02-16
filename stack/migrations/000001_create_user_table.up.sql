CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users
(
    id                    UUID      NOT NULL DEFAULT (uuid_generate_v4()) PRIMARY KEY,
    email                 VARCHAR(256) NOT NULL UNIQUE,
    password                VARCHAR(256) NOT NULL,
    first_name            VARCHAR(256),
    last_name             VARCHAR(256),
    full_name             VARCHAR(256) GENERATED ALWAYS AS (first_name || ' ' || last_name) STORED,
    avatar_url            VARCHAR(256),
    location              VARCHAR(256),
    refresh_token_version INTEGER   NOT NULL DEFAULT 1,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);