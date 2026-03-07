-- Fix casino_bets schema: add missing columns from data contract change
-- The original migration was modified in-place instead of creating a new one

-- Add bet_choice column if not exists
ALTER TABLE casino_bets ADD COLUMN IF NOT EXISTS bet_choice VARCHAR(50) NOT NULL DEFAULT '';

-- Add profit column if not exists
ALTER TABLE casino_bets ADD COLUMN IF NOT EXISTS profit INT NOT NULL DEFAULT 0;

-- Change result from JSONB to VARCHAR(50) if it's still JSONB
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'casino_bets' AND column_name = 'result' AND data_type = 'jsonb'
    ) THEN
        ALTER TABLE casino_bets ALTER COLUMN result DROP DEFAULT;
        ALTER TABLE casino_bets ALTER COLUMN result TYPE VARCHAR(50) USING result::text;
        ALTER TABLE casino_bets ALTER COLUMN result SET DEFAULT '';
    END IF;
END $$;
