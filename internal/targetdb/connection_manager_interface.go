package targetdb

type DatabaseType string

type TargetConnectionManager interface {
	GetConnectionManagerInstance(datasourceName string) (TargetConnectionManager, error)
	LoadDataCollection(databaseType string, dataCollection map[string]any) error
}
