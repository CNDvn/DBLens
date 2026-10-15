package db

import (
	"context"
)

func GetConnectionManager(ctx context.Context, driverName string, datasourceName string) (SourceConnectionManager, error) {
	dbMap := map[string]SourceConnectionManager{
		"mysql": &MysqlConnectionManager{},
		"maria": &MariaConnectionManager{},
	}

	dbInstance := dbMap[driverName]
	if err := dbInstance.newConnectionManager(ctx, datasourceName); err != nil {
		return nil, err
	}

	return dbInstance, nil
}
