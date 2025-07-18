## go-modules

A Go utility library providing generic collection functions and persistent storage capabilities.

### Environment

- Docker
- Make
- Golang

---

### Features

- **Collections**: Generic utility functions for working with slices and other collections
- **Registry**: Generic interface for persistent storage with BadgerDB implementation
- **Type Safety**: Full use of Go generics for type-safe operations
- **Comprehensive Testing**: Extensive test coverage with parallel execution

---

### Usage

```go
import (
    "github.com/viqueen/go-modules/pkg/collections"
    "github.com/viqueen/go-modules/pkg/registry"
)

// Collections example
numbers := []int{1, 2, 3, 4, 5}
strings := collections.Map(numbers, strconv.Itoa)

// Registry example
reg, err := registry.NewBadgerRegistry[MyType]("./data")
item := registry.Item[MyType]{ID: "key", Data: myData}
stored, err := reg.CreateOrUpdate(item)
```

---

### Housekeeping

- lint it and fix it

```bash
make lint-fix
```

- build it

```bash
make build
```