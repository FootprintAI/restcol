-- 0003_drop_redundant_query_indexes.down.sql
--
-- Deliberately empty, as with 0002.
--
-- 0003 drops indexes that duplicate ones AutoMigrate maintains. Recreating them
-- would restore the duplication this migration exists to remove, so a truthful
-- revert does nothing. The indexes that matter were never dropped.

SELECT 1;
