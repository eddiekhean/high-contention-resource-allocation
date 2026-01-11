CREATE TABLE IF NOT EXISTS images (
    id BIGSERIAL PRIMARY KEY,
    s3_key TEXT NOT NULL UNIQUE,
    dhash TEXT NOT NULL,
    dhash_prefix TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_images_dhash_prefix
ON images (dhash_prefix);
