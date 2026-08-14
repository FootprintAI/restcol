-- 0003_drop_redundant_query_indexes.up.sql
--
-- Removes the indexes 0001_add_query_indexes created, and deletes that
-- migration. They duplicate indexes GORM's AutoMigrate already creates from
-- struct tags on the models:
--
--   ModelDocument.ModelProjectID     index:docScope,priority:1
--   ModelDocument.ModelCollectionID  index:docScope,priority:2
--   ModelDocument.CreatedAt          index:docScope,priority:3
--   ModelCollection.ModelProjectID   index
--   ModelSchema.ModelCollectionID    index
--
-- AutoMigrate emits "docScope"; 0001 emitted "idx_restcol_documents_docscope".
-- Same three columns, different names — so 0001's CREATE INDEX IF NOT EXISTS
-- never matched, and applying it produced a second index over identical
-- columns. Six indexes where three suffice, costing write throughput and disk
-- for no read benefit. Measured, see FootprintAI/grandturks#986.
--
-- Only the migration-created names are dropped here. The AutoMigrate ones are
-- the real indexes and must stay; the two sets are distinguishable because
-- GORM quotes the hyphenated table name (idx_restcol-collections_...) where
-- 0001 used underscores (idx_restcol_collections_...).
--
-- A no-op on any environment that never applied 0001, which on the evidence in
-- grandturks#984 is most of them.

DROP INDEX IF EXISTS idx_restcol_documents_docscope;
DROP INDEX IF EXISTS idx_restcol_collections_model_project_id;
DROP INDEX IF EXISTS idx_restcol_collections_schema_model_collection_id;
