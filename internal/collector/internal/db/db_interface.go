package db

import (
	"context"
)

type SourceConnectionManager interface {
	newConnectionManager(ctx context.Context, datasourceName string) error
	CollectData(ctx context.Context) (map[string]any, error)
}
