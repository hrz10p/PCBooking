CREATE TABLE IF NOT EXISTS users (
     id BIGSERIAL PRIMARY KEY,
     created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    fname VARCHAR(255),
    sname VARCHAR(255),
    email VARCHAR(255) UNIQUE NOT NULL,
    user_role VARCHAR(50),
    password_hash bytea NOT NULL,
    activated bool NOT NULL,
)