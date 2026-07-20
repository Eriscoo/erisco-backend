CREATE TABLE IF NOT EXISTS contact_messages (
    id          SERIAL          PRIMARY KEY,
    name        VARCHAR(100)    DEFAULT NULL,
    email       VARCHAR(255)    NOT NULL,
    subject     VARCHAR(255)    DEFAULT NULL,
    phone       VARCHAR(30)     DEFAULT NULL,
    message     TEXT            DEFAULT NULL,
    created_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);
