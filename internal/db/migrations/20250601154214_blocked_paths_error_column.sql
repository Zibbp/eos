-- +goose Up
-- +goose StatementBegin
ALTER TABLE blocked_paths
ADD COLUMN error_text text NOT NULL DEFAULT '';

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE blocked_paths
DROP COLUMN error_text;

-- +goose StatementEnd