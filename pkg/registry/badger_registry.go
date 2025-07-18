// Package registry provides a generic registry interface for storing and retrieving items.
package registry

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger/v4"
)

type badgerRegistry[T any] struct {
	db *badger.DB
}

// Error constants for registry operations.
const (
	ErrorFailedToOpenDB         = "failed to open database"
	ErrorFailedToCreateOrUpdate = "failed to create or update item"
	ErrorFailedToDeleteItem     = "failed to delete item"
	ErrorFailedToReadItem       = "failed to read item"
	ErrorFailedToListItems      = "failed to list items"
	ErrorFailedToListIDs        = "failed to list ids"
)

// NewBadgerRegistry creates a new Badger registry instance.
func NewBadgerRegistry[T any](path string) (Registry[T], error) {
	opts := badger.DefaultOptions(path).WithLoggingLevel(badger.ERROR)

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrorFailedToOpenDB, err)
	}

	return &badgerRegistry[T]{db: db}, nil
}

func (b *badgerRegistry[T]) CreateOrUpdate(item Item[T]) (*Item[T], error) {
	err := b.db.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(item)
		if err != nil {
			return err //nolint:wrapcheck
		}

		return txn.Set([]byte(item.ID), data)
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrorFailedToCreateOrUpdate, err)
	}

	return &item, nil
}

func (b *badgerRegistry[T]) Read(id string) (*Item[T], error) {
	var item Item[T]

	err := b.db.View(func(txn *badger.Txn) error {
		stored, err := txn.Get([]byte(id))
		if err != nil {
			return err //nolint:wrapcheck
		}

		return stored.Value(func(val []byte) error {
			return json.Unmarshal(val, &item)
		})
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrorFailedToReadItem, err)
	}

	return &item, nil
}

func (b *badgerRegistry[T]) Delete(itemID string) (*Item[T], error) {
	item, err := b.Read(itemID)
	if err != nil {
		return nil, err
	}

	err = b.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(itemID))
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrorFailedToDeleteItem, err)
	}

	return item, nil
}

func (b *badgerRegistry[T]) ListItems(filter Filter[T]) ([]Item[T], error) {
	var items []Item[T]

	err := b.db.View(func(txn *badger.Txn) error {
		options := badger.DefaultIteratorOptions
		options.PrefetchValues = true
		options.PrefetchSize = 50

		it := txn.NewIterator(options)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()

			err := item.Value(func(val []byte) error {
				var stored Item[T]
				if err := json.Unmarshal(val, &stored); err != nil {
					return err //nolint:wrapcheck
				}

				if filter(stored) {
					items = append(items, stored)
				}

				return nil
			})
			if err != nil {
				return err //nolint:wrapcheck
			}
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrorFailedToListItems, err)
	}

	return items, nil
}

func (b *badgerRegistry[T]) ListIDs() ([]string, error) {
	var ids []string

	err := b.db.View(func(txn *badger.Txn) error {
		options := badger.DefaultIteratorOptions
		options.PrefetchValues = false
		options.PrefetchSize = 50

		it := txn.NewIterator(options)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			ids = append(ids, string(item.Key()))
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrorFailedToListIDs, err)
	}

	return ids, nil
}
