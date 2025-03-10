CREATE TABLE IF NOT EXISTS users
(
    id INTEGER PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    passhash BYTEA NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_login ON users (login);

CREATE TABLE IF NOT EXISTS apps
(
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);