# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Linting and Code Quality
```bash
# Run linting (read-only)
make lint

# Run linting with auto-fixes
make lint-fix
```

### Building
```bash
# Build using GoReleaser in Docker
make build
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests for specific package
go test ./pkg/collections/
go test ./pkg/registry/

# Run tests with verbose output
go test -v ./pkg/collections/
go test -v ./pkg/registry/

# Run a specific test function
go test -run TestMapTransformIntegersToStrings ./pkg/collections/
go test -run TestRegistry ./pkg/registry/
```

## Project Architecture

This is a Go utility library providing generic collection functions and persistent storage capabilities. The codebase follows these key patterns:

### Package Structure
- `pkg/collections/` - Core collection utilities using Go generics
- `pkg/registry/` - Generic registry interface with BadgerDB implementation for persistent storage
- Tests are in separate `_test.go` files using the `collections_test` package pattern for black-box testing

### Code Organization
- All functions are generic using Go 1.18+ type parameters
- Comprehensive documentation required for all exported functions and packages
- Parallel testing is standard - all tests must include `t.Parallel()`

### Dependencies
The project has strict dependency controls via `.golangci.yaml`:
- Only Go standard library, project modules, and approved libraries allowed
- Use `github.com/stretchr/testify` for assertions
- Use `github.com/ovechkin-dm/mockio/v2` for mocking when needed
- Use `github.com/dgraph-io/badger/v4` for persistent storage in registry package

## Development Environment

### Docker-First Workflow
This project uses Docker containers for consistent development:
- Linting runs in `golangci/golangci-lint:v2.1.2` container
- Builds use `ghcr.io/viqueen/docker-images-golang:main` container
- The Makefile automatically handles Docker flags for CI vs local development

### Code Quality Standards
- All linters enabled by default in `.golangci.yaml`
- Automatic formatting with `gofmt`, `goimports`, and `golines`
- Function length limits enforced (max 60 lines)
- Parallel test execution required
- Package and function documentation mandatory for exported items

### Testing Patterns
When writing tests:
- Use descriptive function names like `TestMapTransformIntegersToStrings`
- Include `t.Parallel()` in every test function
- Test edge cases: nil inputs, empty slices, complex types
- Use black-box testing with separate test packages
- Cover order preservation and type transformations

**Registry Package Testing:**
- Use generic test suites (`testSuite[T any]`) for type-agnostic testing
- Test all CRUD operations: CreateOrUpdate, Read, Delete, ListItems, ListIDs
- Include edge cases: non-existent items, empty registry, filtered lists
- Test data persists in `testdata/badger/` directory during tests

## Registry Package Architecture

The registry package provides a generic interface for persistent storage:

### Core Components
- `Registry[T any]` - Generic interface for CRUD operations on items of type T
- `Item[T any]` - Generic container with ID and Data fields
- `Filter[T any]` - Function type for filtering items based on custom criteria
- `badgerRegistry[T any]` - BadgerDB implementation of the Registry interface

### Usage Patterns
```go
// Create a registry for any type
registry, err := registry.NewBadgerRegistry[MyType]("path/to/db")

// Store and retrieve items
item := registry.Item[MyType]{ID: "key", Data: myData}
stored, err := registry.CreateOrUpdate(item)
retrieved, err := registry.Read("key")

// List and filter items
all, err := registry.ListItems(registry.AllFilter[MyType]())
filtered, err := registry.ListItems(func(item Item[MyType]) bool {
    return item.Data.SomeField == "value"
})
```

### Error Handling
- Structured error constants (ErrorFailedToOpenDB, ErrorFailedToReadItem, etc.)
- Proper error wrapping with context
- BadgerDB errors are wrapped and categorized

## Go Module Information
- Module: `github.com/viqueen/go-modules`
- Go version: 1.24.2
- Target: Generic utility functions for Go developers