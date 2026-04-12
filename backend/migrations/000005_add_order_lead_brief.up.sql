ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS project_type VARCHAR(80),
    ADD COLUMN IF NOT EXISTS budget_range VARCHAR(80),
    ADD COLUMN IF NOT EXISTS city VARCHAR(120),
    ADD COLUMN IF NOT EXISTS contact_method VARCHAR(40);

CREATE INDEX IF NOT EXISTS idx_orders_project_type ON orders(project_type);
CREATE INDEX IF NOT EXISTS idx_orders_city ON orders(city);
