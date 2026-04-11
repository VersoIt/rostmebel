CREATE TABLE IF NOT EXISTS reviews (
    id          BIGSERIAL PRIMARY KEY,
    project_id  BIGINT REFERENCES projects(id) ON DELETE SET NULL,
    order_id    BIGINT REFERENCES orders(id) ON DELETE CASCADE,
    rating      INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment     TEXT NOT NULL,
    images      JSONB NOT NULL DEFAULT '[]', -- [{url}]
    status      VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending|approved|rejected
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_reviews_project ON reviews(project_id) WHERE status = 'approved';
CREATE INDEX IF NOT EXISTS idx_reviews_status ON reviews(status);
