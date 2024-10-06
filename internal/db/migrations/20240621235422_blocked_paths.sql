-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS blocked_paths (
  id uuid PRIMARY KEY,
  path text NOT NULL UNIQUE,
  error_count integer NOT NULL DEFAULT 1,
  is_blocked boolean NOT NULL DEFAULT false,
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
BEFORE UPDATE ON blocked_paths
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd
