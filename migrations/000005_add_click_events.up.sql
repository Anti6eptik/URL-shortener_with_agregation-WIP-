CREATE TABLE click_events (
    id SERIAL PRIMARY KEY,
    short_link_id INT REFERENCES short_links(id),
    clicked_at TIMESTAMPTZ,
    referer TEXT
)