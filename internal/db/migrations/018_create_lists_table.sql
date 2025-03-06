-- Crear tabla para listas genéricas
CREATE TABLE IF NOT EXISTS lists (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Crear tabla para valores de listas
CREATE TABLE IF NOT EXISTS list_values (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    list_id INTEGER NOT NULL,
    value VARCHAR(100) NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    sort_order INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(list_id, value),
    FOREIGN KEY (list_id) REFERENCES lists(id) ON DELETE CASCADE
);

-- Crear índices
CREATE INDEX IF NOT EXISTS idx_list_values_list_id ON list_values(list_id);

-- Insertar lista de disposición
INSERT OR IGNORE INTO lists (name, description) 
VALUES ('disposicion', 'Disposición de la propiedad');

-- Obtener el ID de la lista de disposición e insertar valores
INSERT OR IGNORE INTO list_values (list_id, value, display_name, sort_order)
SELECT id, 'contrafrente', 'Contrafrente', 1 FROM lists WHERE name = 'disposicion';

INSERT OR IGNORE INTO list_values (list_id, value, display_name, sort_order)
SELECT id, 'frente', 'Frente', 2 FROM lists WHERE name = 'disposicion';

INSERT OR IGNORE INTO list_values (list_id, value, display_name, sort_order)
SELECT id, 'interno', 'Interno', 3 FROM lists WHERE name = 'disposicion';

INSERT OR IGNORE INTO list_values (list_id, value, display_name, sort_order)
SELECT id, 'lateral', 'Lateral', 4 FROM lists WHERE name = 'disposicion';

-- Migrar datos existentes usando CASE para normalizar
UPDATE propiedades
SET disposicion = CASE 
    WHEN disposicion IS NULL THEN NULL
    WHEN LOWER(disposicion) LIKE '%contra%frente%' OR LOWER(disposicion) = 'contrafrente' THEN 'contrafrente'
    WHEN LOWER(disposicion) LIKE '%frente%' THEN 'frente'
    WHEN LOWER(disposicion) LIKE '%intern%' THEN 'interno'
    WHEN LOWER(disposicion) LIKE '%lateral%' THEN 'lateral'
    ELSE LOWER(disposicion)
END
WHERE disposicion IS NOT NULL; 