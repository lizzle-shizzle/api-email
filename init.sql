CREATE TABLE IF NOT EXISTS email_log (
    id serial PRIMARY KEY,
    subject text NOT NULL,
    body text NOT NULL,
    email text NOT NULL,
    created_timestamp timestamp with time zone NOT NULL default now()
);