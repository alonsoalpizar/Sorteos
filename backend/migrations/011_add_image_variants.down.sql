-- Remover campos de variantes de imágenes
ALTER TABLE raffle_images
  DROP COLUMN IF EXISTS url_original,
  DROP COLUMN IF EXISTS url_large,
  DROP COLUMN IF EXISTS url_medium,
  DROP COLUMN IF EXISTS url_thumbnail;
