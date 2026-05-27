CREATE TABLE IF NOT EXISTS project_categories (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    slug       VARCHAR(100) NOT NULL UNIQUE,
    icon       VARCHAR(50),
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS projects (
    id                  BIGSERIAL PRIMARY KEY,
    project_category_id BIGINT REFERENCES project_categories(id) ON DELETE SET NULL,
    name                VARCHAR(255) NOT NULL,
    slug                VARCHAR(255) NOT NULL UNIQUE,
    description         TEXT,
    price               NUMERIC(12, 2) NOT NULL,
    price_old           NUMERIC(12, 2),
    images              JSONB NOT NULL DEFAULT '[]',
    specs               JSONB NOT NULL DEFAULT '{}',
    ai_tags             TEXT,
    status              VARCHAR(20) NOT NULL DEFAULT 'draft',
    views_count         INT NOT NULL DEFAULT 0,
    orders_count        INT NOT NULL DEFAULT 0,
    search_vector       TSVECTOR,
    deleted_at          TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_projects_status
    ON projects(status)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_projects_category
    ON projects(project_category_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_projects_search
    ON projects
    USING GIN(search_vector);

CREATE INDEX IF NOT EXISTS idx_projects_price
    ON projects(price)
    WHERE deleted_at IS NULL AND status = 'published';

CREATE OR REPLACE FUNCTION update_project_search_vector()
RETURNS TRIGGER AS $$
DECLARE
    category_name TEXT := '';
    specs_text    TEXT := '';
BEGIN
    IF NEW.project_category_id IS NOT NULL THEN
        SELECT name
          INTO category_name
          FROM project_categories
         WHERE id = NEW.project_category_id;
    END IF;

    SELECT COALESCE(string_agg(spec.key || ' ' || spec.value, ' '), '')
      INTO specs_text
      FROM jsonb_each_text(COALESCE(NEW.specs, '{}'::jsonb)) AS spec(key, value);

    NEW.search_vector := to_tsvector(
        'russian',
        COALESCE(NEW.name, '') || ' ' ||
        COALESCE(category_name, '') || ' ' ||
        COALESCE(NEW.description, '') || ' ' ||
        COALESCE(specs_text, '') || ' ' ||
        COALESCE(NEW.ai_tags, '')
    );

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_project_search_vector ON projects;
CREATE TRIGGER trg_project_search_vector
BEFORE INSERT OR UPDATE ON projects
FOR EACH ROW
EXECUTE FUNCTION update_project_search_vector();

CREATE OR REPLACE FUNCTION refresh_projects_search_vector_for_category()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE projects
       SET name = name
     WHERE project_category_id = NEW.id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_project_category_search_refresh ON project_categories;
CREATE TRIGGER trg_project_category_search_refresh
AFTER UPDATE OF name, slug ON project_categories
FOR EACH ROW
WHEN (OLD.name IS DISTINCT FROM NEW.name OR OLD.slug IS DISTINCT FROM NEW.slug)
EXECUTE FUNCTION refresh_projects_search_vector_for_category();

CREATE TABLE IF NOT EXISTS orders (
    id             BIGSERIAL PRIMARY KEY,
    project_id     BIGINT REFERENCES projects(id) ON DELETE SET NULL,
    client_name    VARCHAR(100) NOT NULL,
    client_phone   VARCHAR(20) NOT NULL,
    client_email   VARCHAR(255),
    comment        TEXT,
    status         VARCHAR(20) NOT NULL DEFAULT 'new',
    ip_address     INET NOT NULL,
    user_agent     TEXT,
    fingerprint    TEXT,
    project_type   VARCHAR(80),
    budget_range   VARCHAR(80),
    city           VARCHAR(120),
    contact_method VARCHAR(40),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_status
    ON orders(status);

CREATE INDEX IF NOT EXISTS idx_orders_created_at
    ON orders(created_at DESC);

CREATE INDEX IF NOT EXISTS idx_orders_ip
    ON orders(ip_address);

CREATE INDEX IF NOT EXISTS idx_orders_project_type
    ON orders(project_type);

CREATE INDEX IF NOT EXISTS idx_orders_city
    ON orders(city);

CREATE TABLE IF NOT EXISTS admins (
    id            BIGSERIAL PRIMARY KEY,
    username      VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    refresh_token VARCHAR(512),
    last_login_at TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ip_blocks (
    id         BIGSERIAL PRIMARY KEY,
    ip_address INET NOT NULL UNIQUE,
    reason     VARCHAR(50) NOT NULL,
    blocked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS reviews (
    id         BIGSERIAL PRIMARY KEY,
    project_id BIGINT REFERENCES projects(id) ON DELETE SET NULL,
    order_id   BIGINT REFERENCES orders(id) ON DELETE CASCADE,
    rating     INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment    TEXT NOT NULL,
    images     JSONB NOT NULL DEFAULT '[]',
    status     VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_reviews_project
    ON reviews(project_id)
    WHERE status = 'approved';

CREATE INDEX IF NOT EXISTS idx_reviews_status
    ON reviews(status);
