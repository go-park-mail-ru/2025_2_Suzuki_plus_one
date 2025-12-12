-- Extensions (OK in a migration)
CREATE EXTENSION IF NOT EXISTS unaccent;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

--------------------------------------------------
-- MEDIA: FTS column

ALTER TABLE media
ADD COLUMN search_doc tsvector GENERATED ALWAYS AS (
  -- title: most important
  setweight(to_tsvector('english', coalesce(title, '')), 'A') ||

  -- plot_summary: very important
  setweight(to_tsvector('english', coalesce(plot_summary, '')), 'B') ||

  -- description: less important
  setweight(to_tsvector('english', coalesce(description, '')), 'C') ||

  -- country + media_type: least important
  setweight(
    to_tsvector(
      'english',
      coalesce(country, '') || ' ' || coalesce(media_type, '')
    ),
    'D'
  )
) STORED;

CREATE INDEX media_search_doc_gin_idx
  ON media
  USING gin (search_doc);

--------------------------------------------------
-- ACTOR: FTS column

ALTER TABLE actor
ADD COLUMN search_doc tsvector GENERATED ALWAYS AS (
  setweight(to_tsvector('english', coalesce(name, '')), 'A') ||
  setweight(to_tsvector('english', coalesce(bio, '')), 'B')
) STORED;

CREATE INDEX actor_search_doc_gin_idx
  ON actor
  USING gin (search_doc);

--------------------------------------------------
-- Trigram indexes for fuzzy search (case-insensitive via lower)

CREATE INDEX media_title_trgm_gin_idx
  ON media
  USING gin (lower(title) gin_trgm_ops);

CREATE INDEX actor_name_trgm_gin_idx
  ON actor
  USING gin (lower(name) gin_trgm_ops);
