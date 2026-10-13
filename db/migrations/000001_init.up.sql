
BEGIN;
CREATE SCHEMA IF NOT EXISTS main_layer;

CREATE TABLE main_layer.sources (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(255) NOT NULL UNIQUE,
    db_type    TEXT NOT NULL,
    host       VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE main_layer.snapshots (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_id           UUID NOT NULL
                            REFERENCES main_layer.sources(id)
                            ON DELETE CASCADE,
    capture_started_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    capture_finished_at TIMESTAMPTZ,
    status              TEXT NOT NULL DEFAULT 'collecting',
    attempt_count       INTEGER NOT NULL DEFAULT 0,
    last_error          TEXT,
    processed_at        TIMESTAMPTZ,

    CONSTRAINT ck_snapshots_status
        CHECK (
            status IN (
                'collecting',
                'collected',
                'processing',
                'processed',
                'failed'
            )
        ),

    CONSTRAINT ck_snapshots_attempt_count
        CHECK (attempt_count >= 0),

    CONSTRAINT ck_snapshots_capture_times
        CHECK (
            capture_finished_at IS NULL
            OR capture_finished_at >= capture_started_at
        )
);

CREATE INDEX idx_snapshots_source_capture_started
    ON main_layer.snapshots (source_id, capture_started_at DESC);

CREATE INDEX idx_snapshots_status_capture_started
    ON main_layer.snapshots (status, capture_started_at);

COMMENT ON TABLE main_layer.snapshots IS
    'A single collection run representing an observed database state.';

CREATE SCHEMA IF NOT EXISTS raw_layer;

-- Raw records from performance_schema.data_locks
CREATE TABLE raw_layer.lock_records (
    id           BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    snapshot_id  UUID NOT NULL
                     REFERENCES main_layer.snapshots(id)
                     ON DELETE CASCADE,
    payload      JSONB NOT NULL,
    collected_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Raw records from performance_schema.data_lock_waits
CREATE TABLE raw_layer.wait_records (
    id           BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    snapshot_id  UUID NOT NULL
                     REFERENCES main_layer.snapshots(id)
                     ON DELETE CASCADE,
    payload      JSONB NOT NULL,
    collected_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Raw records from information_schema.innodb_trx
CREATE TABLE raw_layer.transaction_records (
    id           BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    snapshot_id  UUID NOT NULL
                     REFERENCES main_layer.snapshots(id)
                     ON DELETE CASCADE,
    payload      JSONB NOT NULL,
    collected_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Raw records from performance_schema.threads
CREATE TABLE raw_layer.thread_records (
    id           BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    snapshot_id  UUID NOT NULL
                     REFERENCES main_layer.snapshots(id)
                     ON DELETE CASCADE,
    payload      JSONB NOT NULL,
    collected_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Raw records from performance_schema.events_statements_current
CREATE TABLE raw_layer.statement_records (
    id           BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    snapshot_id  UUID NOT NULL
                     REFERENCES main_layer.snapshots(id)
                     ON DELETE CASCADE,
    payload      JSONB NOT NULL,
    collected_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Snapshot-based retrieval
CREATE INDEX idx_raw_lock_records_snapshot
    ON raw_layer.lock_records (snapshot_id);

CREATE INDEX idx_raw_wait_records_snapshot
    ON raw_layer.wait_records (snapshot_id);

CREATE INDEX idx_raw_transaction_records_snapshot
    ON raw_layer.transaction_records (snapshot_id);

CREATE INDEX idx_raw_thread_records_snapshot
    ON raw_layer.thread_records (snapshot_id);

CREATE INDEX idx_raw_statement_records_snapshot
    ON raw_layer.statement_records (snapshot_id);

COMMENT ON TABLE raw_layer.lock_records IS
    'Raw rows collected from MySQL performance_schema.data_locks.';

COMMENT ON TABLE raw_layer.wait_records IS
    'Raw rows collected from MySQL performance_schema.data_lock_waits.';

COMMENT ON TABLE raw_layer.transaction_records IS
    'Raw rows collected from MySQL information_schema.innodb_trx.';

COMMENT ON TABLE raw_layer.thread_records IS
    'Raw rows collected from MySQL performance_schema.threads.';

COMMENT ON TABLE raw_layer.statement_records IS
    'Raw rows collected from MySQL performance_schema.events_statements_current.';

COMMIT;
