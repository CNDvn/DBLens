package main

import (
	"context"

	"github.com/cndvn/dblens/internal/collector"
)

func main() {
	mysqlCollectorInstance := collector.Collector{}.GetInstace(context.Background(), collector.CollectorConfig{
		SourceDriverName:     "mysql",
		SourceDatasourceName: "root:123456@tcp(localhost:3309)/test",
		TargetDriverName:     "pgx",
		TargetDatasourceName: "postgres://postgres:123456@localhost:5432/postgres?sslmode=disable",
	})
	mysqlCollectorInstance.Run()
}
