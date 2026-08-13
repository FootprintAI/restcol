-- 0002_purge_soft_deleted_documents.up.sql
--
-- DeleteDocument is now a hard delete (#136). Documents deleted before that
-- change are still in the table as tombstones: deleted_at is set, every read
-- filters them out, and nothing can bring them back — but their `data` jsonb is
-- intact, and for inference sinks that column holds base64 image bytes.
--
-- Those rows are the reason the old behaviour failed the retention question in
-- FootprintAI/grandturks#936. Callers were told the documents were deleted.
-- This makes that true.
--
-- IRREVERSIBLE. There is no down migration that can restore the data, only one
-- that describes what was removed. Take a backup first if the rows have any
-- value to you — by definition of the old API they should not, since nothing
-- could read them.

DELETE FROM "restcol-documents" WHERE deleted_at IS NOT NULL;

-- Collections are deliberately NOT purged here. CollectionCURD.Delete is still
-- a soft delete: a collection row carries an id, a type and a summary, not user
-- payload, and hard-deleting one would fire the OnDelete:SET NULL constraint on
-- restcol-collections-schema and leave orphaned schema rows behind. That is a
-- separate change with its own consequences; see #136.
