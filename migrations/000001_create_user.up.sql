CREATE TABLE IF NOT EXISTS users
(
    id       SERIAL PRIMARY KEY,
    email    VARCHAR(255),
    username VARCHAR(255),
    password VARCHAR(255)
);