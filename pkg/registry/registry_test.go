package registry_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/viqueen/go-modules/pkg/registry"
)

type TestData struct {
	Name string
}

func TestRegistry(t *testing.T) {
	tests := map[string]struct {
		registry registry.Registry[TestData]
	}{}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			createItem := func(id string) registry.Item[TestData] {
				return registry.Item[TestData]{ID: id, Data: TestData{Name: "Test"}}
			}
			createUpdated := func(item registry.Item[TestData]) registry.Item[TestData] {
				item.Data.Name = "Updated"

				return item
			}
			suite := newTestSuite[TestData](
				test.registry,
				createItem,
				createUpdated,
			)
			suite.runAll(t)
		})
	}
}

type testSuite[T any] struct {
	registry      registry.Registry[T]
	createItem    func(id string) registry.Item[T]
	createUpdated func(item registry.Item[T]) registry.Item[T]
}

func newTestSuite[T any](
	reg registry.Registry[T],
	createItem func(id string) registry.Item[T],
	createUpdated func(item registry.Item[T]) registry.Item[T],
) *testSuite[T] {
	return &testSuite[T]{
		registry:      reg,
		createItem:    createItem,
		createUpdated: createUpdated,
	}
}

func (s *testSuite[T]) runAll(t *testing.T) {
	t.Helper()
	t.Run("CreateAndRead", s.testCreateAndRead)
	t.Run("Update", s.testUpdate)
	t.Run("Delete", s.testDelete)
	t.Run("ReadNonExistent", s.testReadNonExistent)
	t.Run("DeleteNonExistent", s.testDeleteNonExistent)
	t.Run("ListItems", s.testListItems)
	t.Run("ListIDs", s.testListIDs)
	t.Run("FilteredList", s.testFilteredList)
	t.Run("EmptyRegistry", s.testEmptyRegistry)
	t.Run("MultipleOperations", s.testMultipleOperations)
}

func (s *testSuite[T]) testCreateAndRead(t *testing.T) {
	t.Helper()

	item := s.createItem("test-1")
	created, err := s.registry.CreateOrUpdate(item)
	require.NoError(t, err, "Failed to create item")
	assert.Equal(t, item.ID, created.ID, "Created item ID mismatch")

	read, err := s.registry.Read(item.ID)
	require.NoError(t, err, "Failed to read item")
	assert.Equal(t, item.ID, read.ID, "Read item ID mismatch")
}

func (s *testSuite[T]) testUpdate(t *testing.T) {
	t.Helper()

	item := s.createItem("test-update")

	_, err := s.registry.CreateOrUpdate(item)
	require.NoError(t, err, "Failed to create item for update")

	updated := s.createUpdated(item)
	updated.ID = item.ID // Ensure same ID

	result, err := s.registry.CreateOrUpdate(updated)
	require.NoError(t, err, "Failed to update item")
	assert.Equal(t, updated.ID, result.ID, "Updated item ID mismatch")
}

func (s *testSuite[T]) testDelete(t *testing.T) {
	t.Helper()

	item := s.createItem("test-delete")

	_, err := s.registry.CreateOrUpdate(item)
	require.NoError(t, err, "Failed to create item for deletion")

	deleted, err := s.registry.Delete(item.ID)
	require.NoError(t, err, "Failed to delete item")
	assert.Equal(t, item.ID, deleted.ID, "Deleted item ID mismatch")

	// Verify item is actually deleted
	_, err = s.registry.Read(item.ID)
	assert.Error(t, err, "Expected error when reading deleted item")
}

func (s *testSuite[T]) testReadNonExistent(t *testing.T) {
	t.Helper()

	_, err := s.registry.Read("non-existent-id")
	assert.Error(t, err, "Expected error when reading non-existent item")
}

func (s *testSuite[T]) testDeleteNonExistent(t *testing.T) {
	t.Helper()

	_, err := s.registry.Delete("non-existent-id")
	assert.Error(t, err, "Expected error when deleting non-existent item")
}

func (s *testSuite[T]) testListItems(t *testing.T) {
	t.Helper()
	// Create multiple items
	items := []registry.Item[T]{
		s.createItem("list-1"),
		s.createItem("list-2"),
		s.createItem("list-3"),
	}

	for _, item := range items {
		_, err := s.registry.CreateOrUpdate(item)
		require.NoError(t, err, "Failed to create item %s", item.ID)
	}

	// List all items
	listed, err := s.registry.ListItems(registry.AllFilter[T]())
	require.NoError(t, err, "Failed to list items")

	// Verify all items are present
	for _, item := range items {
		found := false

		for _, listedItem := range listed {
			if listedItem.ID == item.ID {
				found = true

				break
			}
		}

		assert.True(t, found, "Item %s not found in list", item.ID)
	}
}

func (s *testSuite[T]) testListIDs(t *testing.T) {
	t.Helper()
	// Create multiple items
	items := []registry.Item[T]{
		s.createItem("ids-1"),
		s.createItem("ids-2"),
		s.createItem("ids-3"),
	}

	for _, item := range items {
		_, err := s.registry.CreateOrUpdate(item)
		require.NoError(t, err, "Failed to create item %s", item.ID)
	}

	// List all IDs
	ids, err := s.registry.ListIDs(registry.AllFilter[T]())
	require.NoError(t, err, "Failed to list IDs")

	// Verify all IDs are present
	for _, item := range items {
		found := false

		for _, id := range ids {
			if id == item.ID {
				found = true

				break
			}
		}

		assert.True(t, found, "ID %s not found in list", item.ID)
	}
}

func (s *testSuite[T]) testFilteredList(t *testing.T) {
	t.Helper()
	// Create items with specific pattern
	items := []registry.Item[T]{
		s.createItem("filter-match-1"),
		s.createItem("filter-match-2"),
		s.createItem("other-1"),
		s.createItem("other-2"),
	}

	for _, item := range items {
		_, err := s.registry.CreateOrUpdate(item)
		require.NoError(t, err, "Failed to create item %s", item.ID)
	}

	// Filter for items with "filter-match" prefix
	filter := func(item registry.Item[T]) bool {
		return strings.HasPrefix(item.ID, "filter-match")
	}

	// List filtered items
	filtered, err := s.registry.ListItems(filter)
	require.NoError(t, err, "Failed to list filtered items")

	// Verify correct items are returned
	assert.Len(t, filtered, 2, "Expected 2 filtered items")

	for _, item := range filtered {
		assert.Contains(
			t,
			item.ID,
			"filter-match",
			"Filtered item ID should contain 'filter-match'",
		)
	}

	// Test filtered IDs
	filteredIDs, err := s.registry.ListIDs(filter)
	require.NoError(t, err, "Failed to list filtered IDs")
	assert.Len(t, filteredIDs, 2, "Expected 2 filtered IDs")
}

func (s *testSuite[T]) testEmptyRegistry(t *testing.T) {
	t.Helper()
	// Ensure registry is empty by deleting all items
	all, _ := s.registry.ListItems(registry.AllFilter[T]())
	for _, item := range all {
		s.registry.Delete(item.ID)
	}

	// Test listing empty registry
	items, err := s.registry.ListItems(registry.AllFilter[T]())
	require.NoError(t, err, "Failed to list items from empty registry")
	assert.Empty(t, items, "Expected 0 items in empty registry")

	ids, err := s.registry.ListIDs(registry.AllFilter[T]())
	require.NoError(t, err, "Failed to list IDs from empty registry")
	assert.Empty(t, ids, "Expected 0 IDs in empty registry")
}

func (s *testSuite[T]) testMultipleOperations(t *testing.T) {
	t.Helper()
	// Create items
	item1 := s.createItem("multi-1")
	item2 := s.createItem("multi-2")
	item3 := s.createItem("multi-3")

	_, err := s.registry.CreateOrUpdate(item1)
	require.NoError(t, err, "Failed to create item1")

	_, err = s.registry.CreateOrUpdate(item2)
	require.NoError(t, err, "Failed to create item2")

	_, err = s.registry.CreateOrUpdate(item3)
	require.NoError(t, err, "Failed to create item3")

	// Update item2
	updated2 := s.createUpdated(item2)
	updated2.ID = item2.ID

	_, err = s.registry.CreateOrUpdate(updated2)
	require.NoError(t, err, "Failed to update item2")

	// Delete item1
	_, err = s.registry.Delete(item1.ID)
	require.NoError(t, err, "Failed to delete item1")

	// List remaining items
	items, err := s.registry.ListItems(registry.AllFilter[T]())
	require.NoError(t, err, "Failed to list items after operations")

	// Should have 2 items (item2 and item3)
	foundCount := 0

	for _, item := range items {
		if item.ID == item2.ID || item.ID == item3.ID {
			foundCount++
		}
	}

	assert.Equal(t, 2, foundCount, "Expected 2 items after operations")
}
