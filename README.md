# GORM Facade

[![Go Reference](https://pkg.go.dev/badge/github.com/jgomezhrdz/godbkit.svg)](https://pkg.go.dev/github.com/jgomezhrdz/godbkit)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/jgomezhrdz/godbkit)](https://goreportcard.com/report/github.com/jgomezhrdz/godbkit)

GORM Facade is a powerful abstraction layer that provides a clean and consistent interface for working with GORM in Go. It implements the Criteria pattern to build complex queries in a type-safe manner, making database operations more maintainable and testable.

## Why Use This Package?

- **Simplified GORM Usage**: Provides a cleaner, more intuitive API on top of GORM
- **Criteria Pattern**: Build complex queries using a fluent, type-safe Criteria API
- **Consistent API**: Standardized way to perform common database operations
- **Batch Operations**: Efficient batch processing with built-in worker pools
- **Type Safety**: Compile-time checking of query structures

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
go get github.com/jgomezhrdz/godbkit
```

## Features

- **Query Building**: Fluent API for building complex queries
- **Batch Processing**: Concurrent processing of database operations
- **Worker Pool**: Configurable worker pool for concurrent operations
- **Criteria API**: Type-safe query construction
- **Pagination**: Built-in support for paginated results

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

### Core Components

1. **Criteria Manager**: Builds and manages query criteria
2. **Worker Pool**: Manages concurrent database operations
3. **Facade**: Provides a simplified interface to GORM operations

## Quick Start

```go
package main

import (
	"context"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/jgomezhrdz/godbkit/facade"
	"github.com/jgomezhrdz/godbkit/criteriamanager"
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

### Criteria Manager

The `criteriamanager` package provides a fluent API for building complex queries:

```go
criteria := criteriamanager.NewCriteria()
    .AddFilter("name", "=", "John")
    .AddFilter("age", ">", 21)
    .SetOrder("created_at DESC")
    .SetLimit(10)
    .SetOffset(0)
```

### Worker Pool

The `worker` package provides a worker pool for concurrent processing:

```go
pool := worker.NewPool(5) // Create a pool with 5 workers
pool.Start()

defer pool.Stop()

// Submit tasks to the pool
for _, task := range tasks {
    pool.Submit(task)
}
```

### Facade

The `facade` package provides a simplified interface to GORM operations:

```go
gormDB, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
dbFacade := facade.NewGormFacade(gormDB)

// Use the facade for database operations
err := dbFacade.Create(context.Background(), &user)
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create a new branch for your feature
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [GORM](https://gorm.io/) - The fantastic ORM library for Go
- [Work Pool Pattern](https://gobyexample.com/worker-pools) - For the worker pool implementation
