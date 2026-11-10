CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY,
    email         TEXT        NOT NULL UNIQUE,
    username      TEXT        NOT NULL,

    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_id ON users (id);
