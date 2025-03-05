-- +goose Up
-- +goose StatementBegin
-- Primero, crear una columna temporal para almacenar el ID del tipo de propiedad
ALTER TABLE propiedades ADD COLUMN property_type_id INTEGER;

-- Actualizar la columna property_type_id con el ID correspondiente de property_types
UPDATE propiedades
SET property_type_id = (
    SELECT id FROM property_types 
    WHERE name = tipo_propiedad
)
WHERE tipo_propiedad IS NOT NULL;

-- Crear una columna de respaldo para el tipo de propiedad original
ALTER TABLE propiedades ADD COLUMN tipo_propiedad_original TEXT;

-- Guardar los valores originales
UPDATE propiedades
SET tipo_propiedad_original = tipo_propiedad
WHERE tipo_propiedad IS NOT NULL;

-- Eliminar la columna tipo_propiedad
ALTER TABLE propiedades DROP COLUMN tipo_propiedad;

-- Renombrar property_type_id a tipo_propiedad
ALTER TABLE propiedades RENAME COLUMN property_type_id TO tipo_propiedad;

-- Crear un índice para mejorar el rendimiento de las consultas
CREATE INDEX IF NOT EXISTS idx_propiedades_tipo_propiedad ON propiedades(tipo_propiedad);

-- Agregar restricción de clave foránea
-- Nota: SQLite no permite agregar restricciones de clave foránea después de crear la tabla,
-- por lo que esto es más bien documentativo. La integridad referencial debe ser manejada por la aplicación.
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Crear una columna temporal para restaurar el tipo de propiedad como texto
ALTER TABLE propiedades ADD COLUMN tipo_propiedad_text TEXT;

-- Restaurar los valores originales si están disponibles
UPDATE propiedades
SET tipo_propiedad_text = (
    SELECT name FROM property_types 
    WHERE id = tipo_propiedad
)
WHERE tipo_propiedad IS NOT NULL;

-- Si no hay valores originales disponibles, usar los valores de respaldo
UPDATE propiedades
SET tipo_propiedad_text = tipo_propiedad_original
WHERE tipo_propiedad_text IS NULL AND tipo_propiedad_original IS NOT NULL;

-- Eliminar la columna tipo_propiedad (que ahora es un ID)
ALTER TABLE propiedades DROP COLUMN tipo_propiedad;

-- Eliminar la columna de respaldo
ALTER TABLE propiedades DROP COLUMN tipo_propiedad_original;

-- Renombrar tipo_propiedad_text a tipo_propiedad
ALTER TABLE propiedades RENAME COLUMN tipo_propiedad_text TO tipo_propiedad;

-- Eliminar el índice
DROP INDEX IF EXISTS idx_propiedades_tipo_propiedad;
-- +goose StatementEnd 