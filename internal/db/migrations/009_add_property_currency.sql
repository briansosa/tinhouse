-- +goose Up
-- +goose StatementBegin
ALTER TABLE propiedades ADD COLUMN moneda TEXT DEFAULT 'USD';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE propiedades DROP COLUMN moneda;
-- +goose StatementEnd 