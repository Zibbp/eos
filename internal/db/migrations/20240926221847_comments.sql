-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS comments (
  id text PRIMARY KEY,
  text text NOT NULL,
  timestamp timestamptz NOT NULL,
  like_count integer,
  is_favorited boolean DEFAULT false,
  author text NOT NULL,
  author_id text NOT NULL,
  author_thumbnail text,
  author_is_uploader boolean DEFAULT false,
  parent text,
  video_id uuid NOT NULL,
  FOREIGN KEY (video_id) REFERENCES videos(id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS comments;
-- +goose StatementEnd
