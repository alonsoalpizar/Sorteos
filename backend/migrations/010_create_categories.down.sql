-- Eliminar columna category_id de raffles
ALTER TABLE raffles DROP CONSTRAINT IF EXISTS fk_raffles_category;
DROP INDEX IF EXISTS idx_raffles_category;
ALTER TABLE raffles DROP COLUMN IF EXISTS category_id;

-- Eliminar índices de categories
DROP INDEX IF EXISTS idx_categories_display_order;
DROP INDEX IF EXISTS idx_categories_is_active;
DROP INDEX IF EXISTS idx_categories_slug;

-- Eliminar tabla categories
DROP TABLE IF EXISTS categories;
