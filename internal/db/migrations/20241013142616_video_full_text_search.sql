-- +goose Up
-- +goose StatementBegin
ALTER TABLE videos
ADD COLUMN title_fts_en tsvector
GENERATED ALWAYS AS (to_tsvector('english', title)) STORED;

CREATE INDEX title_fts_en_idx
ON videos USING gin (title_fts_en);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS title_fts_en_idx;
ALTER TABLE videos DROP COLUMN IF EXISTS title_fts_en;
-- +goose StatementEnd
