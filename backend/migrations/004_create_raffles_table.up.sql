-- Create ENUM types for raffle
CREATE TYPE raffle_status AS ENUM ('draft', 'active', 'suspended', 'completed', 'cancelled');

-- Create raffles table
CREATE TABLE raffles (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT uuid_generate_v4(),

    -- Owner
    user_id BIGINT NOT NULL,

    -- Basic info
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status raffle_status DEFAULT 'draft' NOT NULL,

    -- Pricing
    price_per_number DECIMAL(10,2) NOT NULL CHECK (price_per_number > 0),
    total_numbers INT NOT NULL CHECK (total_numbers > 0),
    min_number INT NOT NULL DEFAULT 0,
    max_number INT NOT NULL DEFAULT 99,

    -- Draw info
    draw_date TIMESTAMP WITH TIME ZONE NOT NULL,
    draw_method VARCHAR(50) DEFAULT 'loteria_nacional_cr' NOT NULL,

    -- Winner info
    winner_number VARCHAR(10),
    winner_user_id BIGINT,

    -- Counters
    sold_count INT DEFAULT 0 NOT NULL,
    reserved_count INT DEFAULT 0 NOT NULL,

    -- Revenue
    total_revenue DECIMAL(12,2) DEFAULT 0 NOT NULL,
    platform_fee_percentage DECIMAL(5,2) DEFAULT 10.00 NOT NULL,
    platform_fee_amount DECIMAL(10,2) DEFAULT 0 NOT NULL,
    net_amount DECIMAL(12,2) DEFAULT 0 NOT NULL,

    -- Settlement
    settled_at TIMESTAMP WITH TIME ZONE,
    settlement_status VARCHAR(50) DEFAULT 'pending',

    -- Metadata
    metadata JSONB,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    published_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Foreign keys
    CONSTRAINT fk_raffles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
    CONSTRAINT fk_raffles_winner FOREIGN KEY (winner_user_id) REFERENCES users(id) ON DELETE SET NULL,

    -- Constraints
    CONSTRAINT check_max_greater_than_min CHECK (max_number > min_number),
    CONSTRAINT check_sold_count_valid CHECK (sold_count >= 0 AND sold_count <= total_numbers),
    CONSTRAINT check_reserved_count_valid CHECK (reserved_count >= 0 AND reserved_count <= total_numbers),
    CONSTRAINT check_draw_date_future CHECK (draw_date > created_at)
);

-- Create indexes for performance
CREATE INDEX idx_raffles_user_id ON raffles(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_raffles_status ON raffles(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_raffles_draw_date ON raffles(draw_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_raffles_created_at ON raffles(created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_raffles_active ON raffles(status, draw_date) WHERE status = 'active' AND deleted_at IS NULL;

-- Create trigger for updated_at
CREATE OR REPLACE FUNCTION update_raffles_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_raffles_updated_at
    BEFORE UPDATE ON raffles
    FOR EACH ROW
    EXECUTE FUNCTION update_raffles_updated_at();

-- Create trigger for automatic calculations
CREATE OR REPLACE FUNCTION calculate_raffle_revenue()
RETURNS TRIGGER AS $$
BEGIN
    -- Calculate total revenue
    NEW.total_revenue = NEW.sold_count * NEW.price_per_number;

    -- Calculate platform fee
    NEW.platform_fee_amount = NEW.total_revenue * (NEW.platform_fee_percentage / 100);

    -- Calculate net amount for raffle owner
    NEW.net_amount = NEW.total_revenue - NEW.platform_fee_amount;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_calculate_raffle_revenue
    BEFORE INSERT OR UPDATE OF sold_count, price_per_number, platform_fee_percentage ON raffles
    FOR EACH ROW
    EXECUTE FUNCTION calculate_raffle_revenue();

-- Comments
COMMENT ON TABLE raffles IS 'Tabla principal de sorteos/rifas';
COMMENT ON COLUMN raffles.uuid IS 'UUID público del sorteo';
COMMENT ON COLUMN raffles.draw_method IS 'Método de sorteo: loteria_nacional_cr, manual, random';
COMMENT ON COLUMN raffles.platform_fee_percentage IS 'Porcentaje de comisión de la plataforma (default 10%)';
COMMENT ON COLUMN raffles.settlement_status IS 'Estado de liquidación: pending, processing, completed, failed';
