CREATE TABLE IF NOT EXISTS publics (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    photo text NOT NULL,
    telegraph_link text NOT NULL,
    username text NOT NULL,
    link_to_user text NOT NULL,
    link_to_public text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL default NOW(),
    version integer NOT NULL DEFAULT 1
);