CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    password VARCHAR(50) NOT NULL
);

INSERT INTO users (username, password) VALUES
('admin','admin');

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    category VARCHAR(20) NOT NULL,
    name VARCHAR(50) NOT NULL,
    price INTEGER,
    material VARCHAR(20) NOT NULL,
    brand VARCHAR(20) NOT NULL,
    produce_time VARCHAR(10) NOT NULL,
    image VARCHAR(20) NOT NULL
);