-- +goose Up
-- +goose StatementBegin
ALTER TABLE propiedades ADD COLUMN condicion TEXT;
ALTER TABLE propiedades ADD COLUMN orientacion TEXT;
ALTER TABLE propiedades ADD COLUMN disposicion TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE propiedades DROP COLUMN disposicion;
ALTER TABLE propiedades DROP COLUMN orientacion;
ALTER TABLE propiedades DROP COLUMN condicion;
-- +goose StatementEnd 