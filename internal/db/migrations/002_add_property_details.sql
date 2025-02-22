ALTER TABLE propiedades ADD COLUMN tipo_propiedad VARCHAR(50);
ALTER TABLE propiedades ADD COLUMN ubicacion VARCHAR(100);
ALTER TABLE propiedades ADD COLUMN dormitorios INTEGER;
ALTER TABLE propiedades ADD COLUMN banios INTEGER;
ALTER TABLE propiedades ADD COLUMN antiguedad VARCHAR(50);
ALTER TABLE propiedades ADD COLUMN superficie_cubierta FLOAT;
ALTER TABLE propiedades ADD COLUMN superficie_total FLOAT;
ALTER TABLE propiedades ADD COLUMN frente FLOAT;
ALTER TABLE propiedades ADD COLUMN fondo FLOAT;
ALTER TABLE propiedades ADD COLUMN ambientes INTEGER;
ALTER TABLE propiedades ADD COLUMN expensas FLOAT;
ALTER TABLE propiedades ADD COLUMN descripcion TEXT;
ALTER TABLE propiedades ADD COLUMN status VARCHAR(20) DEFAULT 'pending'; 