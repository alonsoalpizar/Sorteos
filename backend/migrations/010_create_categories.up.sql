-- Crear tabla de categorías
CREATE TABLE IF NOT EXISTS categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    icon VARCHAR(50) NOT NULL,  -- emoji
    description TEXT,
    display_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Índices
CREATE INDEX idx_categories_slug ON categories(slug);
CREATE INDEX idx_categories_is_active ON categories(is_active);
CREATE INDEX idx_categories_display_order ON categories(display_order);

-- Insertar categorías iniciales
INSERT INTO categories (name, slug, icon, description, display_order, is_active) VALUES
('Electrónica', 'electronica', '📱', 'Smartphones, laptops, tablets, consolas y más', 1, true),
('Vehículos', 'vehiculos', '🚗', 'Carros, motos, bicicletas y accesorios', 2, true),
('Hogar', 'hogar', '🏠', 'Electrodomésticos, muebles y decoración', 3, true),
('Otros', 'otros', '🎁', 'Premios variados y especiales', 4, true);

-- Agregar columna category_id a raffles
ALTER TABLE raffles ADD COLUMN IF NOT EXISTS category_id BIGINT;
ALTER TABLE raffles ADD CONSTRAINT fk_raffles_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL;
CREATE INDEX idx_raffles_category ON raffles(category_id);

-- Asignar categoría por defecto "Otros" a sorteos existentes
UPDATE raffles
SET category_id = (SELECT id FROM categories WHERE slug = 'otros')
WHERE category_id IS NULL;

-- Comentarios
COMMENT ON TABLE categories IS 'Categorías de sorteos (electrónica, vehículos, hogar, etc)';
COMMENT ON COLUMN categories.slug IS 'Slug único para URLs amigables';
COMMENT ON COLUMN categories.icon IS 'Emoji representativo de la categoría';
COMMENT ON COLUMN categories.display_order IS 'Orden de visualización en filtros';
COMMENT ON COLUMN raffles.category_id IS 'Categoría a la que pertenece el sorteo';
