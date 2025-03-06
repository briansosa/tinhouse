-- +goose Up
-- +goose StatementBegin

-- 1. Crear una columna temporal para almacenar el valor numérico de antiguedad
ALTER TABLE propiedades ADD COLUMN antiguedad_int INTEGER;

-- 2. Crear una columna de respaldo para el valor original de antiguedad
ALTER TABLE propiedades ADD COLUMN antiguedad_original TEXT;

-- 3. Guardar los valores originales
UPDATE propiedades
SET antiguedad_original = antiguedad
WHERE antiguedad IS NOT NULL;

-- 4. Eliminar la columna antiguedad
ALTER TABLE propiedades DROP COLUMN antiguedad;

-- 5. Renombrar la columna antiguedad_int a antiguedad
ALTER TABLE propiedades RENAME COLUMN antiguedad_int TO antiguedad;

-- 6. Crear índice para el nuevo campo
CREATE INDEX IF NOT EXISTS idx_propiedades_antiguedad ON propiedades(antiguedad);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- 1. Crear una columna temporal para restaurar el valor original
ALTER TABLE propiedades ADD COLUMN antiguedad_text TEXT;

-- 2. Restaurar los valores originales
UPDATE propiedades
SET antiguedad_text = antiguedad_original
WHERE antiguedad_original IS NOT NULL;

-- 3. Eliminar la columna antiguedad
ALTER TABLE propiedades DROP COLUMN antiguedad;

-- 4. Renombrar la columna antiguedad_text a antiguedad
ALTER TABLE propiedades RENAME COLUMN antiguedad_text TO antiguedad;

-- 5. Eliminar la columna de respaldo
ALTER TABLE propiedades DROP COLUMN antiguedad_original;

-- 6. Eliminar el índice
DROP INDEX IF EXISTS idx_propiedades_antiguedad;

-- +goose StatementEnd 