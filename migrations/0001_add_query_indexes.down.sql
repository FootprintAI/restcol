-- 0001_add_query_indexes.down.sql
-- Reverts 0001_add_query_indexes.up.sql.

DROP INDEX IF EXISTS idx_restcol_documents_docscope;
DROP INDEX IF EXISTS idx_restcol_collections_schema_model_collection_id;
DROP INDEX IF EXISTS idx_restcol_collections_model_project_id;
