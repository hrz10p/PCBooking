CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    fname VARCHAR(255),
    sname VARCHAR(255),
    email VARCHAR(255) UNIQUE NOT NULL,
    user_role VARCHAR(50),
    password_hash BYTEA NOT NULL,
    activated BOOL NOT NULL
);
