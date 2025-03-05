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
    antiguedad TEXT,
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

-- Índices para mejorar performance
CREATE INDEX IF NOT EXISTS idx_propiedades_codigo ON propiedades(codigo);
CREATE INDEX IF NOT EXISTS idx_propiedades_inmobiliaria ON propiedades(inmobiliaria_id);
CREATE INDEX IF NOT EXISTS idx_busquedas_fecha ON busquedas(created_at);
CREATE INDEX IF NOT EXISTS idx_property_ratings_property_id ON property_ratings(property_id);
CREATE INDEX IF NOT EXISTS idx_property_notes_property_id ON property_notes(property_id); 