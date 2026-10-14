package targetdb

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PgxConnectionManager struct {
	connection *sql.DB
}

var datasourceInstance *PgxConnectionManager

func (t *PgxConnectionManager) GetConnectionManagerInstance(datasourceName string) (TargetConnectionManager, error) {
	if datasourceInstance == nil {

		db, err := sql.Open("pgx", datasourceName)
		if err != nil {
			return nil, err
		}

		datasourceInstance = &PgxConnectionManager{
			connection: db,
		}
	}
	return datasourceInstance, nil
}

func (t *PgxConnectionManager) InsertLockRecords(snapshotId string, payload string) error {
	query := `
		insert into raw_layer.lock_records (snapshot_id, payload, collected_at)
		values ($1, $2, $3)
	`
	_, err := t.connection.Exec(query, snapshotId, payload, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (t *PgxConnectionManager) InsertStatementRecords(snapshotId string, payload string) error {
	query := `
		insert into raw_layer.statement_records (snapshot_id, payload, collected_at)
		values ($1, $2, $3)
	`
	_, err := t.connection.Exec(query, snapshotId, payload, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (t *PgxConnectionManager) InsertThreadRecords(snapshotId string, payload string) error {
	query := `
		insert into raw_layer.thread_records (snapshot_id, payload, collected_at)
		values ($1, $2, $3)
	`
	_, err := t.connection.Exec(query, snapshotId, payload, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (t *PgxConnectionManager) InsertTransactionRecords(snapshotId string, payload string) error {
	query := `
		insert into raw_layer.transaction_records (snapshot_id, payload, collected_at)
		values ($1, $2, $3)
	`
	_, err := t.connection.Exec(query, snapshotId, payload, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (t *PgxConnectionManager) InsertWaitsRecords(snapshotId string, payload string) error {
	query := `
		insert into raw_layer.wait_records (snapshot_id, payload, collected_at)
		values ($1, $2, $3)
	`
	_, err := t.connection.Exec(query, snapshotId, payload, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (t *PgxConnectionManager) CreateSource(host string, name string, dbType string) (*SourceEntity, error) {
	query := `
		INSERT INTO main_layer.sources (id, host, name, db_type)
		VALUES ($1, $2, $3, $4)
		RETURNING id, host, name, db_type, created_at
	`
	id := uuid.New().String()
	_, err := t.connection.Exec(query, id, host, name, dbType)
	if err != nil {
		return nil, err
	}
	return &SourceEntity{
		Id:        id,
		Host:      host,
		Name:      name,
		DbType:    dbType,
		CreatedAt: time.Now(),
	}, nil
}

func (t *PgxConnectionManager) GetSourceByName(name string) (*SourceEntity, error) {
	query := `
		SELECT id, host, name, db_type, created_at
		FROM main_layer.sources
		WHERE name = $1
		`

	sourceEntity := &SourceEntity{}
	err := t.connection.QueryRow(query, name).Scan(&sourceEntity.Id, &sourceEntity.Host, &sourceEntity.Name, &sourceEntity.DbType, &sourceEntity.CreatedAt)
	if err != nil {
		return nil, err
	}
	return sourceEntity, nil
}

func (t *PgxConnectionManager) InitSnapshot(sourceId string) (*SnapshotEntity, error) {
	query := `
		INSERT INTO main_layer.snapshots (id, source_id, capture_started_at,  status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, source_id, capture_started_at, status
	`
	id := uuid.New().String()

	now := time.Now()
	snapshotEntity := &SnapshotEntity{
		Id:               id,
		SourceId:         sourceId,
		CaptureStartedAt: &now,
		Status:           "collecting",
	}

	_, err := t.connection.Exec(query,
		snapshotEntity.Id,
		snapshotEntity.SourceId,
		*snapshotEntity.CaptureStartedAt,
		snapshotEntity.Status,
	)
	if err != nil {
		return nil, err
	}
	return snapshotEntity, nil
}

func (t *PgxConnectionManager) UpdateStatusSnapshot(snapshotId string, status string) error {
	query := `
		UPDATE main_layer.snapshots
		SET status = $1
		WHERE id = $2
	`
	_, err := t.connection.Exec(query, status, snapshotId)
	if err != nil {
		return err
	}
	return nil
}
