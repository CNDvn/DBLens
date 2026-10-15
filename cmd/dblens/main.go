package main

import (
	"context"
	"os"
	"sync"

	"github.com/cndvn/dblens/internal/collector"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	godotenv.Load()
	sourceDatasourceName := os.Getenv("SOURCE_DATASOURCE_NAME")
	targetDatasourceName := os.Getenv("TARGET_DATASOURCE_NAME")

	// mysqlCollectorInstance := collector.Collector{}.GetInstace(context.Background(), collector.CollectorConfig{
	// 	SourceDriverName:     "mysql",
	// 	SourceDatasourceName: sourceDatasourceName,
	// 	SourceHost:           "192.168.1.1",
	// 	SourceName:           "dblens",
	// 	TargetDriverName:     "pgx",
	// 	TargetDatasourceName: targetDatasourceName,
	// })

	mariaCollectorInstance := collector.Collector{}.GetInstace(context.Background(), collector.CollectorConfig{
		SourceDriverName:     "maria",
		SourceDatasourceName: sourceDatasourceName,
		SourceHost:           "vn production",
		SourceName:           "prod vn",
		TargetDriverName:     "pgx",
		TargetDatasourceName: targetDatasourceName,
	})

	wg := sync.WaitGroup{}
	wg.Add(1)

	c := cron.New(cron.WithSeconds())

	_, err := c.AddFunc("@every 2m", func() {
		mariaCollectorInstance.Run()
	})
	if err != nil {
		panic(err)
	}

	c.Start()

	wg.Wait()
}
