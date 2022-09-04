CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    tg_id BIGINT UNIQUE,
    name VARCHAR(255),
    is_admin BOOLEAN,
    is_banned BOOLEAN
);

CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE,
    owner_id BIGINT,
    owner_tag VARCHAR(255),
    members VARCHAR(255),
    date_created DATE,
    deleted BOOLEAN
);

CREATE TABLE IF NOT EXISTS matches (
    date_created  DATE PRIMARY KEY,
    matches_queue bytea
);

CREATE TABLE IF NOT EXISTS tokens (
    id SERIAL PRIMARY KEY,
    token VARCHAR(255),
    used BOOLEAN
);