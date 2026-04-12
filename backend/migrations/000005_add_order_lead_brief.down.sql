DROP INDEX IF EXISTS idx_orders_city;
DROP INDEX IF EXISTS idx_orders_project_type;

ALTER TABLE orders
    DROP COLUMN IF EXISTS contact_method,
    DROP COLUMN IF EXISTS city,
    DROP COLUMN IF EXISTS budget_range,
    DROP COLUMN IF EXISTS project_type;
