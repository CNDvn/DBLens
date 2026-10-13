BEGIN;

DROP TABLE IF EXISTS raw_layer.statement_records;
DROP TABLE IF EXISTS raw_layer.thread_records;
DROP TABLE IF EXISTS raw_layer.transaction_records;
DROP TABLE IF EXISTS raw_layer.wait_records;
DROP TABLE IF EXISTS raw_layer.lock_records;

DROP TABLE IF EXISTS main_layer.snapshots;
DROP TABLE IF EXISTS main_layer.sources;

COMMIT;


