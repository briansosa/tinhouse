-- Up
ALTER TABLE propiedades ADD COLUMN imagenes TEXT;

-- Down
ALTER TABLE propiedades DROP COLUMN imagenes; 