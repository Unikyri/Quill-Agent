DROP INDEX IF EXISTS uq_ingestion_jobs_universe_hash;
ALTER TABLE ingestion_jobs DROP COLUMN IF EXISTS content_hash;
