-- Up
-- Primero eliminamos la tabla existente
DROP TABLE IF EXISTS property_ratings;

-- Recreamos la tabla con la restricción UNIQUE
CREATE TABLE property_ratings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    property_id INTEGER NOT NULL,
    rating TEXT NOT NULL CHECK(rating IN ('like', 'dislike')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (property_id) REFERENCES propiedades(id),
    UNIQUE(property_id)
);

-- Recreamos el índice
CREATE INDEX idx_property_ratings_property_id ON property_ratings(property_id);

-- Down
DROP INDEX IF EXISTS idx_property_ratings_property_id;
DROP TABLE IF EXISTS property_ratings;
CREATE TABLE property_ratings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    property_id INTEGER NOT NULL,
    rating TEXT NOT NULL CHECK(rating IN ('like', 'dislike')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (property_id) REFERENCES propiedades(id)
);
CREATE INDEX idx_property_ratings_property_id ON property_ratings(property_id); 