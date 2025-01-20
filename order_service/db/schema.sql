CREATE TABLE "orders" (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    note VARCHAR,
    status BIGINT NOT NULL CHECK (status IN (1, 2, 3))
);
