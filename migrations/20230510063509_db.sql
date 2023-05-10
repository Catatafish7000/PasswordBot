-- +goose Up
-- +goose StatementBegin
CREATE TABLE passwords(
    id int,
    login text,
    password text,
    created_at timestamp,
    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
