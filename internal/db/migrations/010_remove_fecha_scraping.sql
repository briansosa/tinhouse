-- +goose Up
-- +goose StatementBegin
ALTER TABLE propiedades DROP COLUMN fecha_scraping;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE propiedades ADD COLUMN fecha_scraping TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
-- +goose StatementEnd 