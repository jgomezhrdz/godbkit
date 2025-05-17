package godbkit

import (
	"context"
	"godbkit/criteriamanager"
)

type godbkit interface {
	Select(ctx context.Context, criteria criteriamanager.Criteria, model interface{}, response interface{}, selectQuery string, joinQuery []string, conditions []string) error
	CreateBatch(ctx context.Context, table string, data interface{}, batchSize *int) error
	Create(ctx context.Context, model interface{}) error
	Update(ctx context.Context, idField string, model interface{}) error
	UpdateBatch(ctx context.Context, idField string, data interface{}, workerSize *int) error
}
