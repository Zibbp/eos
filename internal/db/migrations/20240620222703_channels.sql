-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS channels (
  id uuid PRIMARY KEY,
  ext_id text NOT NULL UNIQUE,
  name text NOT NULL UNIQUE,
  description text,
  image_path text,
  generate_thumbnails boolean NOT NULL DEFAULT false,
  created_at timestamptz NOT NULL DEFAULT current_timestamp,
  updated_at timestamptz NOT NULL DEFAULT current_timestamp
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_updated_at
BEFORE UPDATE ON channels
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd
