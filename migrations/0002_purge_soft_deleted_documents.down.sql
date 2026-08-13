-- 0002_purge_soft_deleted_documents.down.sql
--
-- Deliberately empty. 0002 deletes rows; nothing here can bring them back.
--
-- A down migration that silently does nothing is better than one that pretends:
-- restoring these documents requires a backup taken before 0002 ran, not SQL.
-- Reverting the code change (making DeleteDocument scoped again) does not
-- resurrect them either.

SELECT 1;
