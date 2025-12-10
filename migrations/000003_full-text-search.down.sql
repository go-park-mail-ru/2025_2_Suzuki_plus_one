-- Down migration for full-text search + trigram features.
-- Reverse changes made in the corresponding up migration.

-- Drop trigram GIN indexes
DROP INDEX IF EXISTS actor_name_trgm_gin_idx;
DROP INDEX IF EXISTS media_title_trgm_gin_idx;

-- Drop full-text GIN indexes
DROP INDEX IF EXISTS actor_search_doc_gin_idx;
DROP INDEX IF EXISTS media_search_doc_gin_idx;

-- Drop generated search_doc columns
ALTER TABLE actor DROP COLUMN IF EXISTS search_doc;
ALTER TABLE media DROP COLUMN IF EXISTS search_doc;

-- Optionally drop extensions that were created in the up migration.
-- Only drop if you are sure no other code relies on them.
DROP EXTENSION IF EXISTS pg_trgm;
DROP EXTENSION IF EXISTS unaccent;