-- +goose Up
-- +goose StatementBegin
-- Crear tabla de tipos de propiedad
CREATE TABLE IF NOT EXISTS property_types (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insertar tipos de propiedad predeterminados
INSERT OR IGNORE INTO property_types (code, name) VALUES 
    ('house', 'Casa'),
    ('apartment', 'Departamento'),
    ('ph', 'PH'),
    ('local', 'Local'),
    ('office', 'Oficina'),
    ('land', 'Terreno'),
    ('warehouse', 'Galpón');

-- Crear índice para búsquedas rápidas por código
CREATE INDEX IF NOT EXISTS idx_property_types_code ON property_types(code);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_property_types_code;
DROP TABLE IF EXISTS property_types;
-- +goose StatementEnd 