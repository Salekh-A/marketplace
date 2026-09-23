CREATE TABLE IF NOT EXISTS analytics_orders (
                                                id UUID PRIMARY KEY,
                                                product_id UUID NOT NULL,
                                                quantity INTEGER NOT NULL CHECK (quantity > 0),
    price NUMERIC(12, 2) NOT NULL CHECK (price >= 0),
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL
    );

CREATE INDEX IF NOT EXISTS idx_analytics_orders_product_id
    ON analytics_orders(product_id);

CREATE INDEX IF NOT EXISTS idx_analytics_orders_status
    ON analytics_orders(status);