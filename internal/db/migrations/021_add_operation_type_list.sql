-- +goose Up
-- +goose StatementBegin
-- Insertar la lista de tipos de operación
INSERT INTO lists (name, description, created_at, updated_at)
VALUES ('tipo_operacion', 'Tipos de operación de propiedades', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Obtener el ID de la lista recién creada e insertar los valores estandarizados
INSERT INTO list_values (list_id, value, display_name, sort_order, created_at, updated_at)
VALUES 
    ((SELECT id FROM lists WHERE name = 'tipo_operacion'), 'Venta', 'Venta', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ((SELECT id FROM lists WHERE name = 'tipo_operacion'), 'Alquiler', 'Alquiler', 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ((SELECT id FROM lists WHERE name = 'tipo_operacion'), 'Alquiler Temporario', 'Alquiler Temporario', 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Eliminar los valores de la lista
DELETE FROM list_values WHERE list_id = (SELECT id FROM lists WHERE name = 'tipo_operacion');

-- Eliminar la lista
DELETE FROM lists WHERE name = 'tipo_operacion';
-- +goose StatementEnd 