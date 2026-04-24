-- 0001_add_query_indexes.up.sql
-- Adds the indexes needed to keep ListByProjectID, GetLatestSchema,
-- ModelDocument.Query, CountByCollection, and DeleteByCollection off of
-- sequential scans as data volume grows.
--
-- Applying this migration is only required for environments that run with
-- --restcol_auto_migrate=false (production). Dev environments using
-- AutoMigrate pick these up from the gorm struct tags on next boot.
--
-- All statements are idempotent (IF NOT EXISTS) so it is safe to re-run.

-- ModelCollection.ModelProjectID: backs ListByProjectID.
CREATE INDEX IF NOT EXISTS idx_restcol_collections_model_project_id
    ON "restcol-collections" (model_project_id);

-- ModelSchema.ModelCollectionID: backs GetLatestSchema.
CREATE INDEX IF NOT EXISTS idx_restcol_collections_schema_model_collection_id
    ON "restcol-collections-schema" (model_collection_id);

-- ModelDocument composite (project, collection, created_at): backs Query,
-- CountByCollection, and DeleteByCollection. Leading columns support
-- equality filters; trailing created_at supports the ORDER BY in Query.
CREATE INDEX IF NOT EXISTS idx_restcol_documents_docscope
    ON "restcol-documents" (model_project_id, model_collection_id, created_at);
