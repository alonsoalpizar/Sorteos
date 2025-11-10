-- Create raffle_images table
CREATE TABLE raffle_images (
    id BIGSERIAL PRIMARY KEY,

    -- Raffle reference
    raffle_id BIGINT NOT NULL,

    -- Image info
    filename VARCHAR(255) NOT NULL,
    original_filename VARCHAR(255) NOT NULL,
    file_path VARCHAR(512) NOT NULL,
    file_size BIGINT NOT NULL CHECK (file_size > 0),
    mime_type VARCHAR(100) NOT NULL,

    -- Image metadata
    width INT,
    height INT,
    alt_text VARCHAR(255),

    -- Ordering
    display_order INT DEFAULT 0 NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE NOT NULL,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Foreign keys
    CONSTRAINT fk_raffle_images_raffle FOREIGN KEY (raffle_id) REFERENCES raffles(id) ON DELETE CASCADE,

    -- Constraints
    CONSTRAINT check_valid_mime_type CHECK (
        mime_type IN ('image/jpeg', 'image/jpg', 'image/png', 'image/webp', 'image/gif')
    ),
    CONSTRAINT check_max_file_size CHECK (file_size <= 10485760) -- 10 MB max
);

-- Create indexes
CREATE INDEX idx_raffle_images_raffle_id ON raffle_images(raffle_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_raffle_images_primary ON raffle_images(raffle_id) WHERE is_primary = TRUE AND deleted_at IS NULL;
CREATE INDEX idx_raffle_images_order ON raffle_images(raffle_id, display_order) WHERE deleted_at IS NULL;

-- Create unique partial index to ensure only one primary image per raffle
CREATE UNIQUE INDEX uq_one_primary_per_raffle
    ON raffle_images(raffle_id)
    WHERE is_primary = TRUE AND deleted_at IS NULL;

-- Create trigger for updated_at
CREATE OR REPLACE FUNCTION update_raffle_images_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_raffle_images_updated_at
    BEFORE UPDATE ON raffle_images
    FOR EACH ROW
    EXECUTE FUNCTION update_raffle_images_updated_at();

-- Create trigger to automatically set primary image if it's the first one
CREATE OR REPLACE FUNCTION auto_set_primary_image()
RETURNS TRIGGER AS $$
DECLARE
    v_image_count INT;
BEGIN
    -- Count existing images for this raffle
    SELECT COUNT(*) INTO v_image_count
    FROM raffle_images
    WHERE raffle_id = NEW.raffle_id
      AND deleted_at IS NULL;

    -- If this is the first image, make it primary
    IF v_image_count = 0 THEN
        NEW.is_primary = TRUE;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_auto_set_primary_image
    BEFORE INSERT ON raffle_images
    FOR EACH ROW
    EXECUTE FUNCTION auto_set_primary_image();

-- Create trigger to prevent deleting the primary image if there are other images
CREATE OR REPLACE FUNCTION prevent_delete_only_primary()
RETURNS TRIGGER AS $$
DECLARE
    v_other_images_count INT;
BEGIN
    -- Only check if the image being deleted is primary
    IF OLD.is_primary = TRUE AND NEW.deleted_at IS NOT NULL THEN
        -- Count other non-deleted images
        SELECT COUNT(*) INTO v_other_images_count
        FROM raffle_images
        WHERE raffle_id = OLD.raffle_id
          AND id != OLD.id
          AND deleted_at IS NULL;

        -- If there are other images, we need to set one as primary
        IF v_other_images_count > 0 THEN
            -- Set the next image as primary (by display_order, then by id)
            UPDATE raffle_images
            SET is_primary = TRUE
            WHERE id = (
                SELECT id
                FROM raffle_images
                WHERE raffle_id = OLD.raffle_id
                  AND id != OLD.id
                  AND deleted_at IS NULL
                ORDER BY display_order ASC, id ASC
                LIMIT 1
            );
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_prevent_delete_only_primary
    BEFORE UPDATE OF deleted_at ON raffle_images
    FOR EACH ROW
    WHEN (OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL)
    EXECUTE FUNCTION prevent_delete_only_primary();

-- Comments
COMMENT ON TABLE raffle_images IS 'Imágenes de los sorteos (galería)';
COMMENT ON COLUMN raffle_images.filename IS 'Nombre del archivo almacenado (ej: uuid.jpg)';
COMMENT ON COLUMN raffle_images.original_filename IS 'Nombre original del archivo subido por el usuario';
COMMENT ON COLUMN raffle_images.file_path IS 'Ruta completa del archivo (local o S3)';
COMMENT ON COLUMN raffle_images.display_order IS 'Orden de visualización en la galería (0 = primero)';
COMMENT ON COLUMN raffle_images.is_primary IS 'Indica si es la imagen principal del sorteo';
COMMENT ON INDEX uq_one_primary_per_raffle IS 'Asegura que solo haya una imagen primary por sorteo';
