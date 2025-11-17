-- Agregar campos para variantes de imágenes optimizadas
ALTER TABLE raffle_images
  ADD COLUMN IF NOT EXISTS url_original VARCHAR(500),
  ADD COLUMN IF NOT EXISTS url_large VARCHAR(500),
  ADD COLUMN IF NOT EXISTS url_medium VARCHAR(500),
  ADD COLUMN IF NOT EXISTS url_thumbnail VARCHAR(500);

-- Comentarios para documentación
COMMENT ON COLUMN raffle_images.url_original IS 'URL de imagen original (max 1200x1200, JPG/PNG)';
COMMENT ON COLUMN raffle_images.url_large IS 'URL de imagen grande (800x800, WebP 85%)';
COMMENT ON COLUMN raffle_images.url_medium IS 'URL de imagen mediana (400x400, WebP 80%)';
COMMENT ON COLUMN raffle_images.url_thumbnail IS 'URL de thumbnail (150x150, WebP 75%)';
