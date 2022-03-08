CREATE TABLE IF NOT EXISTS ads (
    id bigserial PRIMARY KEY,
    user_id integer,
    link text,
    msg text,
    created_at timestamp(0) with time zone NOT NULL default NOW(),
    paid boolean NOT NULL DEFAULT FALSE,
    version integer NOT NULL DEFAULT 1
);

