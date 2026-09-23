CREATE TABLE IF NOT EXISTS orders (
                                      id UUID PRIMARY KEY,
                                      product_id UUID NOT NULL,
                                      quantity INTEGER NOT NULL CHECK (quantity > 0),
    price NUMERIC(12, 2) NOT NULL CHECK (price >= 0),
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
    );

CREATE INDEX IF NOT EXISTS idx_orders_product_id
    ON orders(product_id);

CREATE INDEX IF NOT EXISTS idx_orders_status
    ON orders(status);