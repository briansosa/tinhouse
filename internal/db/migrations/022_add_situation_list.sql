-- +goose Up
-- +goose StatementBegin
-- Insertar la lista de situaciones
INSERT INTO lists (name, description, created_at, updated_at)
VALUES ('situacion', 'Situaciones de propiedades', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Obtener el ID de la lista recién creada e insertar los valores estandarizados
INSERT INTO list_values (list_id, value, display_name, sort_order, created_at, updated_at)
VALUES 
    ((SELECT id FROM lists WHERE name = 'situacion'), 'Vacía', 'Vacía', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ((SELECT id FROM lists WHERE name = 'situacion'), 'Habitada', 'Habitada', 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Eliminar los valores de la lista
DELETE FROM list_values WHERE list_id = (SELECT id FROM lists WHERE name = 'situacion');

-- Eliminar la lista
DELETE FROM lists WHERE name = 'situacion';
-- +goose StatementEnd 