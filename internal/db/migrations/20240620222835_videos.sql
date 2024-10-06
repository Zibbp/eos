-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS videos (
  id uuid PRIMARY KEY,
  ext_id text UNIQUE,
  title text NOT NULL,
  description text,
  upload_date timestamp with time zone NOT NULL,
  uploader text,
  duration integer NOT NULL,
  view_count bigint NOT NULL,
  like_count bigint,
  dislike_count bigint,
  format text,
  height integer,
  width integer,
  resolution text,
  fps real,
  video_codec text,
  vbr real,
  audio_codec text,
  abr real,
  comment_count integer,
  video_path text NOT NULL,
  thumbnail_path text NOT NULL,
  info_path text NOT NULL,
  subtitle_path text[],
  path text NOT NULL,
  storyboard_path text,
  created_at timestamptz NOT NULL DEFAULT current_timestamp,
  updated_at timestamptz NOT NULL DEFAULT current_timestamp,
  channel_id uuid NOT NULL,
  FOREIGN KEY (channel_id) REFERENCES channels(id)
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_updated_at
BEFORE UPDATE ON videos
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd