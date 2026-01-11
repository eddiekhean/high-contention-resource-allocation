CREATE TABLE IF NOT EXISTS images (
    id BIGSERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    dhash BIGINT NOT NULL,
    dhash_prefix SMALLINT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_images_dhash_prefix
ON images (dhash_prefix);
