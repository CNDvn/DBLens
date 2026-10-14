package targetdb

import "time"

type DatabaseType string

type SourceEntity struct {
	Id        string    `json:"id"`
	Host      string    `json:"host"`
	Name      string    `json:"name"`
	DbType    string    `json:"db_type"`
	CreatedAt time.Time `json:"created_at"`
}

type SnapshotEntity struct {
	Id                string     `json:"id"`
	SourceId          string     `json:"source_id"`
	CaptureStartedAt  *time.Time `json:"capture_started_at"`
	CaptureFinishedAt *time.Time `json:"capture_finished_at"`
	Status            string     `json:"status"`
	AttemptCount      int        `json:"attempt_count"`
	LastError         string     `json:"last_error"`
	ProcessedAt       *time.Time `json:"processed_at"`
}

type LockRecordsEntity struct {
	Id          string            `json:"id"`
	SnapshotId  string            `json:"snapshot_id"`
	Payload     map[string]string `json:"payload"`
	CollectedAt *time.Time        `json:"collected_at"`
}

type StatementRecordsEntity struct {
	Id          string            `json:"id"`
	SnapshotId  string            `json:"snapshot_id"`
	Payload     map[string]string `json:"payload"`
	CollectedAt *time.Time        `json:"collected_at"`
}

type ThreadRecordsEntity struct {
	Id          string            `json:"id"`
	SnapshotId  string            `json:"snapshot_id"`
	Payload     map[string]string `json:"payload"`
	CollectedAt *time.Time        `json:"collected_at"`
}

type TransactionRecordsEntity struct {
	Id          string            `json:"id"`
	SnapshotId  string            `json:"snapshot_id"`
	Payload     map[string]string `json:"payload"`
	CollectedAt *time.Time        `json:"collected_at"`
}

type WaitRecordsEntity struct {
	Id          string            `json:"id"`
	SnapshotId  string            `json:"snapshot_id"`
	Payload     map[string]string `json:"payload"`
	CollectedAt *time.Time        `json:"collected_at"`
}

type TargetConnectionManager interface {
	GetConnectionManagerInstance(datasourceName string) (TargetConnectionManager, error)

	CreateSource(host string, name string, dbType string) (*SourceEntity, error)
	GetSourceByName(name string) (*SourceEntity, error)

	InitSnapshot(sourceId string) (*SnapshotEntity, error)
	UpdateStatusSnapshot(snapshotId string, status string) error

	InsertLockRecords(snapshotId string, payload string) error
	InsertStatementRecords(snapshotId string, payload string) error
	InsertThreadRecords(snapshotId string, payload string) error
	InsertTransactionRecords(snapshotId string, payload string) error
	InsertWaitsRecords(snapshotId string, payload string) error
}
