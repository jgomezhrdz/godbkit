package godbkit

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/jgomezhrdz/godbkit/criteriamanager"
	worker "github.com/jgomezhrdz/godbkit/internal/worker"

	"gorm.io/gorm"
)

// Select executes a query with filters, joins, ordering, limits and conditions.
func Select(
	ctx context.Context,
	db *gorm.DB,
	criteria criteriamanager.Criteria,
	model interface{},
	response interface{},
	selectQuery string,
	joinQuery []string,
	conditions []string,
) error {
	query, values := criteriamanager.ParseConditions(criteria.GETFILTROS())

	queryBuilder := db.WithContext(ctx).
		Model(model).
		Select(selectQuery)

	for _, join := range joinQuery {
		queryBuilder = queryBuilder.Joins(join)
	}

	if limit := criteria.GETLIMIT(); limit != nil {
		queryBuilder = queryBuilder.Limit(*limit)
	}
	if offset := criteria.GETOFFSET(); offset != nil {
		queryBuilder = queryBuilder.Offset(*offset)
	}
	if order := criteria.GETORDER(); order != "" {
		queryBuilder = queryBuilder.Order(order)
	}

	if len(conditions) > 0 {
		if query != "" {
			query = "(" + query + ") AND "
		}
		query += "(" + strings.Join(conditions, " AND ") + ")"
	}

	return queryBuilder.Where(query, values...).Scan(response).Error
}

// CreateBatch creates multiple records in batches.
func CreateBatch[Records any](db *gorm.DB, ctx context.Context, table string, data []Records, batchSize *int) error {
	size := 5 // default batch size
	if batchSize != nil {
		size = *batchSize
	}
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return fmt.Errorf("data must be a slice or array")
	}
	if val.Len() == 0 {
		return nil
	}

	tx := db.WithContext(ctx).Table(table).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := tx.CreateInBatches(data, size).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// Create creates a single record and updates NULL for nil pointer fields.
func Create(db *gorm.DB, ctx context.Context, model interface{}) error {
	tx := db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := tx.Create(model).Error; err != nil {
		tx.Rollback()
		return err
	}

	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	t := val.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := val.Field(i)

		if field.Type.Kind() == reflect.Ptr && fieldValue.IsNil() {
			if err := tx.Model(model).Update(field.Name, gorm.Expr("NULL")).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

// Update updates a record by idField, handling NULL for nil pointer fields.
func Update[Records any](db *gorm.DB, ctx context.Context, idField string, model *Records) error {
	tx := db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	err := update(tx, ctx, idField, model)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error trying to modify proyecto on database: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error trying to commit changes on database: %w", err)
	}
	return nil
}

// UpdateBatch updates multiple records concurrently using worker pool.
// Each record is updated in its own transaction.
func UpdateBatch[Records any](db *gorm.DB, ctx context.Context, idField string, data []Records, workerSize *int) error {
	workers := 5 // default batch size
	if workerSize != nil {
		workers = *workerSize
	}
	numRecords := len(data)
	if numRecords == 0 {
		return nil
	}

	pool := worker.NewPool(workers, numRecords)
	pool.Start()

	for i := range data {
		record := &data[i]

		err := pool.Submit(func(args ...any) error {
			tx := db.WithContext(ctx).Begin()
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
					panic(r)
				}
			}()

			err := update(tx, ctx, idField, record)
			if err != nil {
				tx.Rollback()
				return err
			}

			if err := tx.Commit().Error; err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			pool.Stop()
			break
		}
	}

	errs := pool.Wait()
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// update performs the actual update logic on the model.
func update(tx *gorm.DB, ctx context.Context, idField string, model any) error {
	val := reflect.ValueOf(model)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("model must be a pointer")
	}

	elem := val.Elem()
	t := elem.Type()

	idValue := elem.FieldByName(idField)
	if !idValue.IsValid() {
		return fmt.Errorf("field %s does not exist in model", idField)
	}

	query := tx.WithContext(ctx).Model(model).Where(fmt.Sprintf("%s = ?", idField), idValue.Interface())

	if err := query.Updates(model).Error; err != nil {
		return err
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := elem.Field(i)

		if field.Type.Kind() == reflect.Ptr && fieldValue.IsNil() {
			if err := query.Update(field.Name, gorm.Expr("NULL")).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
