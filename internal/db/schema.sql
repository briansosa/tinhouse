-- Tabla de inmobiliarias
CREATE TABLE IF NOT EXISTS inmobiliarias (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT NOT NULL,
    url TEXT,
    sistema TEXT,
    zona TEXT,
    rating REAL,
    direccion TEXT,
    telefono TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de tipos de propiedad
CREATE TABLE IF NOT EXISTS property_types (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code TEXT NOT NULL UNIQUE,  -- Código interno (ej: 'house', 'apartment')
    name TEXT NOT NULL,         -- Nombre para mostrar (ej: 'Casa', 'Departamento')
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insertar tipos de propiedad predeterminados
INSERT OR IGNORE INTO property_types (code, name) VALUES 
    ('house', 'Casa'),
    ('apartment', 'Departamento'),
    ('ph', 'PH');

-- Tabla de búsquedas
CREATE TABLE IF NOT EXISTS busquedas (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    operation TEXT NOT NULL,  -- venta, alquiler
    property_type TEXT NOT NULL,  -- casa, departamento, etc
    zone TEXT NOT NULL,
    location TEXT NOT NULL,
    min_price_usd REAL,
    max_price_usd REAL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de propiedades
CREATE TABLE IF NOT EXISTS propiedades (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    inmobiliaria_id INTEGER,
    codigo TEXT NOT NULL UNIQUE,
    titulo TEXT,
    precio TEXT,
    direccion TEXT,
    url TEXT,
    imagen_url TEXT,
    imagenes TEXT,  -- Nueva columna para array de imágenes
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    tipo_propiedad TEXT,
    ubicacion TEXT,
    dormitorios INTEGER,
    banios INTEGER,
    antiguedad INTEGER,
    superficie_cubierta FLOAT,
    superficie_total FLOAT,
    superficie_terreno FLOAT,
    frente FLOAT,
    fondo FLOAT,
    ambientes INTEGER,
    plantas INTEGER,
    cocheras INTEGER,
    situacion TEXT,
    expensas FLOAT,
    descripcion TEXT,
    status TEXT DEFAULT 'pending',
    operacion TEXT,
    condicion TEXT,
    orientacion TEXT,
    disposicion TEXT,
    FOREIGN KEY (inmobiliaria_id) REFERENCES inmobiliarias(id)
);

-- Tabla intermedia busquedas_propiedades
CREATE TABLE IF NOT EXISTS busquedas_propiedades (
    busqueda_id INTEGER,
    propiedad_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (busqueda_id, propiedad_id),
    FOREIGN KEY (busqueda_id) REFERENCES busquedas(id),
    FOREIGN KEY (propiedad_id) REFERENCES propiedades(id)
);

-- Tabla de calificaciones de propiedades
CREATE TABLE IF NOT EXISTS property_ratings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    property_id INTEGER NOT NULL,
    rating TEXT NOT NULL CHECK(rating IN ('like', 'dislike')),
    is_favorite BOOLEAN DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (property_id) REFERENCES propiedades(id),
    UNIQUE(property_id)
);

-- Tabla de notas de propiedades
CREATE TABLE IF NOT EXISTS property_notes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    property_id INTEGER NOT NULL,
    note TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (property_id) REFERENCES propiedades(id)
);

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

-- Índices para mejorar performance
CREATE INDEX IF NOT EXISTS idx_propiedades_codigo ON propiedades(codigo);
CREATE INDEX IF NOT EXISTS idx_propiedades_inmobiliaria ON propiedades(inmobiliaria_id);
CREATE INDEX IF NOT EXISTS idx_busquedas_fecha ON busquedas(created_at);
CREATE INDEX IF NOT EXISTS idx_property_ratings_property_id ON property_ratings(property_id);
CREATE INDEX IF NOT EXISTS idx_property_notes_property_id ON property_notes(property_id);
CREATE INDEX IF NOT EXISTS idx_property_feature_relations_property_id ON property_feature_relations(property_id);
CREATE INDEX IF NOT EXISTS idx_property_feature_relations_feature_id ON property_feature_relations(feature_id);
CREATE INDEX IF NOT EXISTS idx_property_features_category ON property_features(category); 