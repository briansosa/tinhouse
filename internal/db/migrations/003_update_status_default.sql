-- Actualizar el valor por defecto de la columna status
ALTER TABLE propiedades ALTER COLUMN status SET DEFAULT 'pending';

-- Actualizar los registros existentes que tengan status 'completed' a 'pending'
UPDATE propiedades SET status = 'pending' WHERE status = 'completed'; 