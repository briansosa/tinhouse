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
    fecha_scraping TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
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

-- Índices para mejorar performance
CREATE INDEX IF NOT EXISTS idx_propiedades_codigo ON propiedades(codigo);
CREATE INDEX IF NOT EXISTS idx_propiedades_inmobiliaria ON propiedades(inmobiliaria_id);
CREATE INDEX IF NOT EXISTS idx_busquedas_fecha ON busquedas(created_at); 