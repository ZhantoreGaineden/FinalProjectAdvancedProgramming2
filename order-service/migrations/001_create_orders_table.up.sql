CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    user_email VARCHAR(150) NOT NULL,
    total_price NUMERIC(10,2) NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'created',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    pet_id UUID NOT NULL,
    price NUMERIC(10,2) NOT NULL
);
