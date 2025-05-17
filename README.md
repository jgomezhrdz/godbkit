# Database ORM Facade

[![Go Reference](https://pkg.go.dev/badge/github.com/yourusername/orm-facade.svg)](https://pkg.go.dev/github.com/yourusername/orm-facade)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/orm-facade)](https://goreportcard.com/report/github.com/yourusername/orm-facade)

Database ORM Facade is a powerful abstraction layer that provides a unified interface for working with multiple ORMs in Go. It implements the Criteria pattern to build complex queries in a clean, type-safe manner, making it easy to switch between different database backends while maintaining a consistent API.

## Why Use This Package?

- **ORM Agnostic**: Write database-agnostic code that works with multiple ORMs (currently supports GORM with extensibility for others)
- **Criteria Pattern**: Build complex queries using a fluent, type-safe Criteria API
- **Consistent API**: Same interface across different database backends
- **Transaction Management**: Simplified transaction handling across different ORMs
- **Batch Operations**: Efficient batch processing with worker pools

## Features

- **Multi-ORM Support**: Currently supports GORM with an extensible architecture for adding more ORMs
- **Criteria Pattern**: Build complex queries using a fluent, type-safe Criteria API
- **Unified CRUD Operations**: Consistent interface for Create, Read, Update, and Delete operations across ORMs
- **Advanced Querying**: Support for complex conditions, joins, and result mapping
- **Batch Processing**: Efficient batch operations with configurable batch sizes
- **Concurrent Operations**: Worker pool implementation for high-performance batch processing
- **Transaction Management**: Consistent transaction handling across different ORMs
- **Type Safety**: Compile-time checking of query structures
- **NULL Handling**: Proper handling of NULL values and pointer fields

## Installation

```bash
go get github.com/yourusername/orm-facade
```

## Supported ORMs

- [x] GORM (primary support)
- [ ] SQLBoiler (planned)
- [ ] Ent (planned)
- [ ] XORM (planned)

## Core Concepts

### Criteria Pattern

The Criteria pattern allows you to build complex queries in a fluent, type-safe manner:

```go
criteria := criteriamanager.NewCriteria()
    .AddFilter("name", "LIKE", "%john%")
    .AddFilter("age", ">", 21)
    .SetOrder("created_at DESC")
    .SetLimit(10)
    .SetOffset(0)
```

### ORM Adapters

Each supported ORM has its own adapter that implements the common interface:

```go
type DatabaseFacade interface {
    Select(ctx context.Context, criteria Criteria, model interface{}, result interface{}) error
    Create(ctx context.Context, model interface{}) error
    Update(ctx context.Context, model interface{}) error
    Delete(ctx context.Context, model interface{}) error
    // ... other common operations
}
```

## Quick Start

### Using with GORM

```go
package main

import (
	"context"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/yourusername/orm-facade/facade"
	"github.com/yourusername/orm-facade/criteriamanager"
)

type User struct {
	gorm.Model
	Name  string
	Email *string
}

func main() {
	// Initialize GORM DB
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Auto migrate your schema
	db.AutoMigrate(&User{})

	// Initialize the GORM adapter
	gormFacade := facade.NewGormFacade(db)

	// Create a new user
	email := "user@example.com"
	user := &User{Name: "John Doe", Email: &email}

	err = gormFacade.Create(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	// Query using criteria
	criteria := criteriamanager.NewCriteria().
		AddFilter("name", "LIKE", "%John%").
		SetOrder("created_at DESC")

	var results []User
	err = gormFacade.Select(context.Background(), criteria, &User{}, &results)
	if err != nil {
		log.Fatal(err)
	}

	// Batch operations with worker pool
	users := []User{
		{Name: "Jane Smith"},
		{Name: "Bob Johnson", Email: &email},
	}

	// Process batch with 2 concurrent workers
	err = gormFacade.ProcessInBatches(context.Background(), users, 2, func(user *User) error {
		// Process each user in a separate goroutine
		return gormFacade.Update(context.Background(), user)
	})

	if err != nil {
		log.Fatal(err)
	}
}
```

## Documentation

### Available Methods

#### Create
```go
func Create(db *gorm.DB, ctx context.Context, model interface{}) error
```
Creates a new record and handles NULL values for nil pointer fields.

#### CreateBatch
```go
func CreateBatch[Records any](db *gorm.DB, ctx context.Context, table string, data []Records, batchSize *int) error
```
Creates multiple records in batches with the specified batch size.

#### Update
```go
func Update[Records any](db *gorm.DB, ctx context.Context, idField string, model *Records) error
```
Updates a record by ID and handles NULL values for nil pointer fields.

#### UpdateBatch
```go
func UpdateBatch[Records any](db *gorm.DB, ctx context.Context, idField string, data []Records, workerSize *int) error
```
Updates multiple records concurrently using a worker pool.

#### Select
```go
func Select(
    ctx context.Context,
    db *gorm.DB,
    criteria criteriamanager.Criteria,
    model interface{},
    response interface{},
    selectQuery string,
    joinQuery []string,
    conditions []string,
) error
```
Executes a query with filters, joins, ordering, limits, and conditions.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [GORM](https://gorm.io/) - The fantastic ORM library for Go
