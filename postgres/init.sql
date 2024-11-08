CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(20) UNIQUE NOT NULL,
    password CHAR(60) NOT NULL
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    category VARCHAR(20) NOT NULL,
    name VARCHAR(50) NOT NULL,
    price INTEGER,
    material VARCHAR(20) NOT NULL,
    brand VARCHAR(20) NOT NULL,
    produce_time VARCHAR(10) NOT NULL,
    image VARCHAR(20) UNIQUE NOT NULL
);