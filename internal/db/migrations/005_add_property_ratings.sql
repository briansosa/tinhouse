-- Migración para agregar la tabla de calificaciones de propiedades

-- Up
CREATE TABLE property_ratings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    property_id INTEGER NOT NULL,
    rating TEXT NOT NULL CHECK(rating IN ('like', 'dislike')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (property_id) REFERENCES propiedades(id)
);

-- Índice para optimizar búsquedas por property_id
CREATE INDEX idx_property_ratings_property_id ON property_ratings(property_id);

-- Down
DROP INDEX IF EXISTS idx_property_ratings_property_id;
DROP TABLE IF EXISTS property_ratings; 