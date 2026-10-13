package targetdb

import (
	"database/sql"
	"fmt"

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

func (t *PgxConnectionManager) LoadDataCollection(databaseType string, dataCollection map[string]any) error {
	if databaseType == "mysql" {
		// Load MySQL data into PostgreSQL
		fmt.Println(dataCollection["payloadDataLocks"])
	}

	return nil
}
