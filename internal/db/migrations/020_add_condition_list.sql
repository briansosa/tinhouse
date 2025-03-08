-- +goose Up
-- +goose StatementBegin
-- Insertar la lista de condiciones
INSERT INTO lists (name, description, created_at, updated_at)
VALUES ('condicion', 'Condiciones de propiedades', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Obtener el ID de la lista recién creada e insertar los valores estandarizados
INSERT INTO list_values (list_id, value, display_name, sort_order, created_at, updated_at)
VALUES 
    ((SELECT id FROM lists WHERE name = 'condicion'), 'Excelente', 'Excelente', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ((SELECT id FROM lists WHERE name = 'condicion'), 'Muy bueno', 'Muy bueno', 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ((SELECT id FROM lists WHERE name = 'condicion'), 'Bueno', 'Bueno', 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ((SELECT id FROM lists WHERE name = 'condicion'), 'Regular', 'Regular', 4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ((SELECT id FROM lists WHERE name = 'condicion'), 'A refaccionar', 'A refaccionar', 5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ((SELECT id FROM lists WHERE name = 'condicion'), 'Reciclado', 'Reciclado', 6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Eliminar los valores de la lista
DELETE FROM list_values WHERE list_id = (SELECT id FROM lists WHERE name = 'condicion');

-- Eliminar la lista
DELETE FROM lists WHERE name = 'condicion';
-- +goose StatementEnd 