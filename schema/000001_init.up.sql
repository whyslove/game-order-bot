CREATE TABLE IF NOT EXISTS admins (
    id SERIAL PRIMARY KEY,
    tg_id BIGINT,
    name VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    owner_id BIGINT,
    owner_tag VARCHAR(255),
    date_created DATE,
    deleted BOOLEAN
);

CREATE TABLE IF NOT EXISTS matches (
    date_created  DATE PRIMARY KEY,
    matches_queue bytea
);

CREATE TABLE IF NOT EXISTS tokens (
    token VARCHAR(255),
    used BOOLEAN
);