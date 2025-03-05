-- +goose Up
-- +goose StatementBegin
ALTER TABLE propiedades ADD COLUMN latitud REAL;
ALTER TABLE propiedades ADD COLUMN longitud REAL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE propiedades DROP COLUMN latitud;
ALTER TABLE propiedades DROP COLUMN longitud;
-- +goose StatementEnd 