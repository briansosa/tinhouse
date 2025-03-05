-- +goose Up
-- +goose StatementBegin
ALTER TABLE propiedades ADD COLUMN operacion TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE propiedades DROP COLUMN operacion;
-- +goose StatementEnd 