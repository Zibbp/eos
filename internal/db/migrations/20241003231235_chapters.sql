-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS chapters (
  id uuid PRIMARY KEY,
  title text NOT NULL,
  start_time integer NOT NULL,
  end_time integer NOT NULL,
  video_id uuid NOT NULL,
  FOREIGN KEY (video_id) REFERENCES videos(id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS chapters;
-- +goose StatementEnd
