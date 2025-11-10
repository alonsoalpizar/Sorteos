-- Create ENUM type for raffle number status
CREATE TYPE raffle_number_status AS ENUM ('available', 'reserved', 'sold');

-- Create raffle_numbers table
CREATE TABLE raffle_numbers (
    id BIGSERIAL PRIMARY KEY,

    -- Raffle reference
    raffle_id BIGINT NOT NULL,

    -- Number info
    number VARCHAR(10) NOT NULL,
    status raffle_number_status DEFAULT 'available' NOT NULL,

    -- Buyer info (only when sold)
    user_id BIGINT,
    reservation_id BIGINT,
    payment_id BIGINT,

    -- Reservation tracking
    reserved_at TIMESTAMP WITH TIME ZONE,
    reserved_until TIMESTAMP WITH TIME ZONE,
    reserved_by BIGINT, -- User who reserved (may differ from buyer)

    -- Sale tracking
    sold_at TIMESTAMP WITH TIME ZONE,
    price DECIMAL(10,2),

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,

    -- Foreign keys
    CONSTRAINT fk_raffle_numbers_raffle FOREIGN KEY (raffle_id) REFERENCES raffles(id) ON DELETE CASCADE,
    CONSTRAINT fk_raffle_numbers_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT fk_raffle_numbers_reserved_by FOREIGN KEY (reserved_by) REFERENCES users(id) ON DELETE SET NULL,

    -- Unique constraint: cada número es único por sorteo
    CONSTRAINT uq_raffle_number UNIQUE (raffle_id, number),

    -- Business constraints
    CONSTRAINT check_reserved_until_future CHECK (
        reserved_until IS NULL OR reserved_until > reserved_at
    ),
    CONSTRAINT check_user_id_when_sold CHECK (
        (status = 'sold' AND user_id IS NOT NULL) OR status != 'sold'
    ),
    CONSTRAINT check_price_when_sold CHECK (
        (status = 'sold' AND price IS NOT NULL AND price > 0) OR status != 'sold'
    )
);

-- Create indexes for performance
CREATE INDEX idx_raffle_numbers_raffle_id ON raffle_numbers(raffle_id);
CREATE INDEX idx_raffle_numbers_status ON raffle_numbers(raffle_id, status);
CREATE INDEX idx_raffle_numbers_user_id ON raffle_numbers(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_raffle_numbers_reserved_until ON raffle_numbers(reserved_until) WHERE status = 'reserved';
CREATE INDEX idx_raffle_numbers_available ON raffle_numbers(raffle_id) WHERE status = 'available';

-- Create trigger for updated_at
CREATE OR REPLACE FUNCTION update_raffle_numbers_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_raffle_numbers_updated_at
    BEFORE UPDATE ON raffle_numbers
    FOR EACH ROW
    EXECUTE FUNCTION update_raffle_numbers_updated_at();

-- Create trigger to update raffle counters
CREATE OR REPLACE FUNCTION update_raffle_counters()
RETURNS TRIGGER AS $$
DECLARE
    v_sold_count INT;
    v_reserved_count INT;
BEGIN
    -- Contar vendidos y reservados para el sorteo
    SELECT
        COUNT(*) FILTER (WHERE status = 'sold'),
        COUNT(*) FILTER (WHERE status = 'reserved')
    INTO v_sold_count, v_reserved_count
    FROM raffle_numbers
    WHERE raffle_id = COALESCE(NEW.raffle_id, OLD.raffle_id);

    -- Actualizar contadores en raffles
    UPDATE raffles
    SET
        sold_count = v_sold_count,
        reserved_count = v_reserved_count,
        updated_at = NOW()
    WHERE id = COALESCE(NEW.raffle_id, OLD.raffle_id);

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_raffle_counters
    AFTER INSERT OR UPDATE OF status ON raffle_numbers
    FOR EACH ROW
    EXECUTE FUNCTION update_raffle_counters();

-- Function to release expired reservations (will be called by cron job)
CREATE OR REPLACE FUNCTION release_expired_reservations()
RETURNS TABLE(raffle_id BIGINT, numbers_released INT) AS $$
DECLARE
    v_raffle_record RECORD;
    v_count INT;
BEGIN
    -- Update all expired reservations
    FOR v_raffle_record IN
        SELECT DISTINCT rn.raffle_id, COUNT(*) as count
        FROM raffle_numbers rn
        WHERE rn.status = 'reserved'
          AND rn.reserved_until < NOW()
        GROUP BY rn.raffle_id
    LOOP
        -- Reset to available
        UPDATE raffle_numbers
        SET
            status = 'available',
            reserved_at = NULL,
            reserved_until = NULL,
            reserved_by = NULL,
            reservation_id = NULL,
            updated_at = NOW()
        WHERE raffle_numbers.raffle_id = v_raffle_record.raffle_id
          AND status = 'reserved'
          AND reserved_until < NOW();

        raffle_id := v_raffle_record.raffle_id;
        numbers_released := v_raffle_record.count;
        RETURN NEXT;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- Comments
COMMENT ON TABLE raffle_numbers IS 'Números individuales de cada sorteo con su estado';
COMMENT ON COLUMN raffle_numbers.number IS 'Número del sorteo (formato: 00, 01, 02... o personalizado)';
COMMENT ON COLUMN raffle_numbers.reserved_until IS 'Fecha límite de reserva (típicamente 5 minutos)';
COMMENT ON COLUMN raffle_numbers.reserved_by IS 'Usuario que hizo la reserva (puede ser diferente al comprador final)';
COMMENT ON FUNCTION release_expired_reservations() IS 'Libera reservas expiradas, retorna raffles afectados';
