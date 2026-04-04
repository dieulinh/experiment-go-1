CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    description TEXT,
    user_id INT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW()
);