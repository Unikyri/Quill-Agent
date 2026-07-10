ALTER TABLE ingestion_jobs ADD COLUMN content_hash CHAR(64);
CREATE UNIQUE INDEX uq_ingestion_jobs_universe_hash
    ON ingestion_jobs(universe_id, content_hash)
    WHERE status <> 'failed' AND content_hash IS NOT NULL;
