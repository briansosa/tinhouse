-- Migración para agregar la tabla de notas de propiedades

-- Up
CREATE TABLE property_notes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    property_id INTEGER NOT NULL,
    note TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (property_id) REFERENCES propiedades(id)
);

-- Índice para optimizar búsquedas por property_id
CREATE INDEX idx_property_notes_property_id ON property_notes(property_id);

-- Down
DROP INDEX IF EXISTS idx_property_notes_property_id;
DROP TABLE IF EXISTS property_notes; 