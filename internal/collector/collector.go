package collector

import (
	"context"
	"fmt"
	"log"

	"github.com/cndvn/dblens/internal/collector/internal/db"
	"github.com/cndvn/dblens/internal/targetdb"
)

type CollectorConfig struct {
	SourceDriverName     string
	SourceDatasourceName string
	TargetDriverName     string
	TargetDatasourceName string
}

type Collector struct {
	config                  CollectorConfig
	sourceConnectionManager db.SourceConnectionManager
	targetConnectionManager targetdb.TargetConnectionManager
}

var instances = make(map[string]*Collector)

func (t Collector) GetInstace(ctx context.Context, collectorConfig CollectorConfig) *Collector {
	instanceKey := collectorConfig.SourceDriverName + "_" + collectorConfig.SourceDatasourceName
	instance := instances[instanceKey]
	if instance == nil {
		sourceConnectionManagerInstance, err := db.GetConnectionManager(ctx, collectorConfig.SourceDriverName, collectorConfig.SourceDatasourceName)
		if err != nil {
			log.Fatalf("Failed to create database instance: %v", err)
		}

		targetConnectionManagerInstance, err := (&targetdb.PgxConnectionManager{}).GetConnectionManagerInstance(collectorConfig.TargetDatasourceName)
		if err != nil {
			log.Fatalf("Failed to create target database instance: %v", err)
		}

		instance = &Collector{
			config:                  collectorConfig,
			sourceConnectionManager: sourceConnectionManagerInstance,
			targetConnectionManager: targetConnectionManagerInstance,
		}
		instances[instanceKey] = instance
	}
	return instance
}

func (t *Collector) Run() {
	fmt.Println("Running collector")

	dataCollection, err := t.sourceConnectionManager.CollectData(context.Background())
	if err != nil {
		log.Fatalf("Failed to collect data: %v", err)
	}

	t.targetConnectionManager.LoadDataCollection(t.config.SourceDriverName, dataCollection)
}
