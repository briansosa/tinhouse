-- +goose Up
-- +goose StatementBegin
-- Tabla para almacenar características normalizadas
CREATE TABLE IF NOT EXISTS property_features (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    category TEXT NOT NULL, -- 'servicio', 'ambiente', 'adicional'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de relación entre propiedades y características
CREATE TABLE IF NOT EXISTS property_feature_relations (
    property_id INTEGER NOT NULL,
    feature_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (property_id, feature_id),
    FOREIGN KEY (property_id) REFERENCES propiedades(id) ON DELETE CASCADE,
    FOREIGN KEY (feature_id) REFERENCES property_features(id) ON DELETE CASCADE
);

-- Crear índice para búsquedas rápidas
CREATE INDEX IF NOT EXISTS idx_property_feature_relations_property_id ON property_feature_relations(property_id);
CREATE INDEX IF NOT EXISTS idx_property_feature_relations_feature_id ON property_feature_relations(feature_id);
CREATE INDEX IF NOT EXISTS idx_property_features_category ON property_features(category);

-- Insertar características predefinidas
-- Servicios
INSERT INTO property_features (name, category) VALUES ('Agua Corriente', 'servicio');
INSERT INTO property_features (name, category) VALUES ('Cloaca', 'servicio');
INSERT INTO property_features (name, category) VALUES ('Gas Natural', 'servicio');
INSERT INTO property_features (name, category) VALUES ('Electricidad', 'servicio');
INSERT INTO property_features (name, category) VALUES ('Pavimento', 'servicio');

-- Ambientes
INSERT INTO property_features (name, category) VALUES ('Cocina', 'ambiente');
INSERT INTO property_features (name, category) VALUES ('Comedor diario', 'ambiente');
INSERT INTO property_features (name, category) VALUES ('Lavadero', 'ambiente');
INSERT INTO property_features (name, category) VALUES ('Patio', 'ambiente');
INSERT INTO property_features (name, category) VALUES ('Living', 'ambiente');
INSERT INTO property_features (name, category) VALUES ('Balcón', 'ambiente');
INSERT INTO property_features (name, category) VALUES ('Terraza', 'ambiente');
INSERT INTO property_features (name, category) VALUES ('Jardín', 'ambiente');
INSERT INTO property_features (name, category) VALUES ('Quincho', 'ambiente');
INSERT INTO property_features (name, category) VALUES ('Playroom', 'ambiente');

-- Adicionales
INSERT INTO property_features (name, category) VALUES ('Calefacción', 'adicional');
INSERT INTO property_features (name, category) VALUES ('Apto profesional', 'adicional');
INSERT INTO property_features (name, category) VALUES ('Termo eléctrico', 'adicional');
INSERT INTO property_features (name, category) VALUES ('Luminoso', 'adicional');
INSERT INTO property_features (name, category) VALUES ('Laundry', 'adicional');
INSERT INTO property_features (name, category) VALUES ('Aire acondicionado', 'adicional');
INSERT INTO property_features (name, category) VALUES ('Pileta', 'adicional');
INSERT INTO property_features (name, category) VALUES ('Seguridad', 'adicional');
INSERT INTO property_features (name, category) VALUES ('Amoblado', 'adicional');
INSERT INTO property_features (name, category) VALUES ('Parrilla', 'adicional');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS property_feature_relations;
DROP TABLE IF EXISTS property_features;
-- +goose StatementEnd 