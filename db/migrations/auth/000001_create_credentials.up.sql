CREATE TABLE IF NOT EXISTS credentials
(
    user_id       INTEGER PRIMARY KEY,

    password_hash TEXT        NOT NULL,

    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_credentials_user_id
    ON credentials (user_id);