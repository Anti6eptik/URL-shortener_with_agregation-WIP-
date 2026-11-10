CREATE TABLE short_links (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(16) NOT NULL UNIQUE,
    original_url TEXT NOT NULL,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ
);

CREATE INDEX idx_short_links_user_id ON short_links(user_id);

CREATE INDEX idx_short_links_expires_at ON short_links(expires_at) WHERE expires_at IS NOT NULL;