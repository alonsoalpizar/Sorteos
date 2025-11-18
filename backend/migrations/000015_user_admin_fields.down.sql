ALTER TABLE users
    DROP COLUMN IF EXISTS suspension_reason,
    DROP COLUMN IF EXISTS suspended_by,
    DROP COLUMN IF EXISTS suspended_at,
    DROP COLUMN IF EXISTS last_kyc_review,
    DROP COLUMN IF EXISTS kyc_reviewer;
