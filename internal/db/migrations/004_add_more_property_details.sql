-- Agregar nuevos campos para detalles adicionales de propiedades
ALTER TABLE propiedades ADD COLUMN plantas INTEGER;
ALTER TABLE propiedades ADD COLUMN cocheras INTEGER;
ALTER TABLE propiedades ADD COLUMN situacion VARCHAR(50);
ALTER TABLE propiedades ADD COLUMN superficie_terreno FLOAT; 