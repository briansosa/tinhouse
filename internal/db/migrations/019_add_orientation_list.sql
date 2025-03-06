-- +goose Up
-- +goose StatementBegin
-- Insertar la lista de orientaciones
INSERT INTO lists (name, description, created_at, updated_at)
VALUES ('orientacion', 'Orientaciones de propiedades', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Obtener el ID de la lista recién creada e insertar los valores
INSERT INTO list_values (list_id, value, display_name, sort_order, created_at, updated_at)
VALUES 
    ((SELECT id FROM lists WHERE name = 'orientacion'), 'Norte', 'Norte', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ((SELECT id FROM lists WHERE name = 'orientacion'), 'Noreste', 'Noreste', 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ((SELECT id FROM lists WHERE name = 'orientacion'), 'Este', 'Este', 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ((SELECT id FROM lists WHERE name = 'orientacion'), 'Sudeste', 'Sudeste', 4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ((SELECT id FROM lists WHERE name = 'orientacion'), 'Sur', 'Sur', 5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ((SELECT id FROM lists WHERE name = 'orientacion'), 'Suroeste', 'Suroeste', 6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ((SELECT id FROM lists WHERE name = 'orientacion'), 'Oeste', 'Oeste', 7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ((SELECT id FROM lists WHERE name = 'orientacion'), 'Noroeste', 'Noroeste', 8, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Eliminar los valores de la lista
DELETE FROM list_values WHERE list_id = (SELECT id FROM lists WHERE name = 'orientacion');

-- Eliminar la lista
DELETE FROM lists WHERE name = 'orientacion';
-- +goose StatementEnd 