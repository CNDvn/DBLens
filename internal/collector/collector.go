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
	SourceId             string
	SourceHost           string
	SourceName           string
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
	err := instance.initSource()
	if err != nil {
		log.Fatalf("Failed to initialize source: %v", err)
	}
	return instance
}

func (t *Collector) initSource() error {
	var err error
	sourceEntity, _ := t.targetConnectionManager.GetSourceByName(t.config.SourceName)

	if sourceEntity == nil {
		sourceEntity, err = t.targetConnectionManager.CreateSource(t.config.SourceHost, t.config.SourceName, t.config.SourceDriverName)
		if err != nil {
			return err
		}
	}

	t.config.SourceId = sourceEntity.Id
	t.config.SourceHost = sourceEntity.Host
	t.config.SourceName = sourceEntity.Name
	t.config.SourceDriverName = sourceEntity.DbType
	return nil
}

func (t *Collector) Run() {
	fmt.Println("Running collector")

	snapshotEntity, err := t.targetConnectionManager.InitSnapshot(t.config.SourceId)
	if err != nil {
		log.Fatalf("Failed to initialize snapshot: %v", err)
	}

	dataCollection, err := t.sourceConnectionManager.CollectData(context.Background())
	if err != nil {
		log.Fatalf("Failed to collect data: %v", err)
	}

	err = t.loadDataCollection(snapshotEntity.Id, dataCollection)
	if err != nil {
		log.Fatalf("Failed to load data collection: %v", err)
	}

	err = t.targetConnectionManager.UpdateStatusSnapshot(snapshotEntity.Id, "collected")
	if err != nil {
		log.Fatalf("Failed to update status snapshot: %v", err)
	}
}

func (t *Collector) loadDataCollection(snapshotId string, payload map[string]string) error {
	if t.config.SourceDriverName == "mysql" {
		err := t.targetConnectionManager.InsertLockRecords(snapshotId, payload["payloadDataLocks"])
		if err != nil {
			return fmt.Errorf("failed to insert lock records: %v", err)
		}

		err = t.targetConnectionManager.InsertStatementRecords(snapshotId, payload["payloadEventsStatement"])
		if err != nil {
			return fmt.Errorf("failed to insert statement records: %v", err)
		}

		err = t.targetConnectionManager.InsertThreadRecords(snapshotId, payload["payloadThreads"])
		if err != nil {
			return fmt.Errorf("failed to insert thread records: %v", err)
		}

		err = t.targetConnectionManager.InsertTransactionRecords(snapshotId, payload["payloadInnoDBTrx"])
		if err != nil {
			return fmt.Errorf("failed to insert transaction records: %v", err)
		}

		err = t.targetConnectionManager.InsertWaitsRecords(snapshotId, payload["payloadDataLockWaits"])
		if err != nil {
			return fmt.Errorf("failed to insert waits records: %v", err)
		}
	}

	return nil
}
