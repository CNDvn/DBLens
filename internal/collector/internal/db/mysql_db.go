package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/cndvn/dblens/pkg/utils"
	_ "github.com/go-sql-driver/mysql"
)

type MysqlConnectionManager struct {
	connection *sql.DB
}

func (t *MysqlConnectionManager) newConnectionManager(ctx context.Context, datasourceName string) error {
	db, err := sql.Open("mysql", datasourceName)
	if err != nil {
		return err
	}
	t.connection = db
	return nil
}

func (t *MysqlConnectionManager) CollectData(ctx context.Context) (map[string]string, error) {
	payloadDataLocks, err := t.collectDataFromDB(ctx, "select * from performance_schema.data_locks")
	if err != nil {
		return nil, err
	}

	payloadDataLockWaits, err := t.collectDataFromDB(ctx, "select * from performance_schema.data_lock_waits")
	if err != nil {
		return nil, err
	}

	payloadInnoDBTrx, err := t.collectDataFromDB(ctx, "select * from information_schema.innodb_trx")
	if err != nil {
		return nil, err
	}

	payloadThreads, err := t.collectDataFromDB(ctx, "select * from performance_schema.threads")
	if err != nil {
		return nil, err
	}

	payloadEventsStatement, err := t.collectDataFromDB(ctx, "select * from performance_schema.events_statements_current")
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"payloadDataLocks":       payloadDataLocks,
		"payloadDataLockWaits":   payloadDataLockWaits,
		"payloadInnoDBTrx":       payloadInnoDBTrx,
		"payloadThreads":         payloadThreads,
		"payloadEventsStatement": payloadEventsStatement,
	}, nil
}

func (t *MysqlConnectionManager) collectDataFromDB(ctx context.Context, query string) (string, error) {
	rows, err := t.connection.QueryContext(ctx, query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	payloads, err := utils.ScanResultDB(rows)
	if err != nil {
		return "", fmt.Errorf("failed to scan result: %v", err)
	}

	jsonData, err := json.Marshal(payloads)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %v", err)
	}

	return string(jsonData), nil
}
