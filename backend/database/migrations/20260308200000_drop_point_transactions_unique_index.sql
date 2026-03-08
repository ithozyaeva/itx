-- Drop unique index that prevents multiple casino bets (and other repeated point transactions)
-- The constraint was attempted to be dropped in 20260306000000 but DROP CONSTRAINT
-- does not remove an index created via CREATE UNIQUE INDEX
DROP INDEX IF EXISTS idx_point_transactions_unique;
