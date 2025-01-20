CREATE TABLE track_user (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    counter INT DEFAULT 0
);