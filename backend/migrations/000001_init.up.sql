-- Categories
CREATE TABLE IF NOT EXISTS categories (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    slug        VARCHAR(100) NOT NULL UNIQUE,
    icon        VARCHAR(50),
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Products
CREATE TABLE IF NOT EXISTS products (
    id              BIGSERIAL PRIMARY KEY,
    category_id     BIGINT REFERENCES categories(id) ON DELETE SET NULL,
    name            VARCHAR(255) NOT NULL,
    slug            VARCHAR(255) NOT NULL UNIQUE,
    description     TEXT,
    price           NUMERIC(12, 2) NOT NULL,
    price_old       NUMERIC(12, 2),
    images          JSONB NOT NULL DEFAULT '[]', -- [{url, is_main}]
    specs           JSONB NOT NULL DEFAULT '{}', -- {"Материал": "Дуб", "Размер": "120x60"}
    ai_tags         TEXT,                        -- "скандинавский, дуб, светлый, спальня"
    status          VARCHAR(20) NOT NULL DEFAULT 'draft', -- published|draft|archived
    views_count     INT NOT NULL DEFAULT 0,
    orders_count    INT NOT NULL DEFAULT 0,
    search_vector   TSVECTOR,                    -- для full-text поиска
    deleted_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_products_status ON products(status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_products_category ON products(category_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_products_search ON products USING GIN(search_vector);
CREATE INDEX IF NOT EXISTS idx_products_price ON products(price) WHERE deleted_at IS NULL AND status = 'published';

-- Trigger for search vector
CREATE OR REPLACE FUNCTION update_product_search_vector()
RETURNS TRIGGER AS $$
BEGIN
    NEW.search_vector := to_tsvector('russian',
        COALESCE(NEW.name, '') || ' ' ||
        COALESCE(NEW.description, '') || ' ' ||
        COALESCE(NEW.ai_tags, '')
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_product_search_vector ON products;
CREATE TRIGGER trg_product_search_vector
BEFORE INSERT OR UPDATE ON products
FOR EACH ROW EXECUTE FUNCTION update_product_search_vector();

-- Orders
CREATE TABLE IF NOT EXISTS orders (
    id              BIGSERIAL PRIMARY KEY,
    product_id      BIGINT REFERENCES products(id) ON DELETE SET NULL,
    client_name     VARCHAR(100) NOT NULL,
    client_phone    VARCHAR(20) NOT NULL,
    client_email    VARCHAR(255),
    comment         TEXT,
    status          VARCHAR(20) NOT NULL DEFAULT 'new', -- new|processing|done|rejected|spam
    ip_address      INET NOT NULL,
    user_agent      TEXT,
    fingerprint     VARCHAR(64),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_orders_ip ON orders(ip_address);

-- Admins
CREATE TABLE IF NOT EXISTS admins (
    id              BIGSERIAL PRIMARY KEY,
    username        VARCHAR(50) NOT NULL UNIQUE,
    password_hash   VARCHAR(255) NOT NULL,
    refresh_token   VARCHAR(512),
    last_login_at   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- IP Blocks
CREATE TABLE IF NOT EXISTS ip_blocks (
    id          BIGSERIAL PRIMARY KEY,
    ip_address  INET NOT NULL UNIQUE,
    reason      VARCHAR(50) NOT NULL, -- spam|rate_limit
    blocked_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at  TIMESTAMPTZ NOT NULL
);
