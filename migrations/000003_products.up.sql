CREATE TABLE products (
    product_id SERIAL PRIMARY KEY,
    product_name VARCHAR(200) UNIQUE NOT NULL,
    quantity VARCHAR(60) NOT NULL
); 