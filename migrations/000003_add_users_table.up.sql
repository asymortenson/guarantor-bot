CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    user_id BIGINT UNIQUE,
    created_at timestamp(0) with time zone NOT NULL default NOW(),
    version integer NOT NULL DEFAULT 1
);

