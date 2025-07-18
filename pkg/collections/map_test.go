package collections_test

import (
	"strconv"
	"testing"

	"github.com/viqueen/go-modules/pkg/collections"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	t.Run("transforms slice of integers to strings", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		expected := []string{"1", "2", "3", "4", "5"}

		result := collections.Map(input, strconv.Itoa)

		assert.Equal(t, expected, result)
	})

	t.Run("transforms slice of strings to their lengths", func(t *testing.T) {
		input := []string{"hello", "world", "go", "test"}
		expected := []int{5, 5, 2, 4}

		result := collections.Map(input, func(s string) int {
			return len(s)
		})

		assert.Equal(t, expected, result)
	})

	t.Run("handles empty slice", func(t *testing.T) {
		var input []int

		var expected []string

		result := collections.Map(input, strconv.Itoa)

		assert.Equal(t, expected, result)
	})

	t.Run("handles nil slice", func(t *testing.T) {
		var input []int

		result := collections.Map(input, strconv.Itoa)

		assert.Nil(t, result)
	})

	t.Run("preserves order", func(t *testing.T) {
		input := []int{3, 1, 4, 1, 5}
		expected := []int{6, 2, 8, 2, 10}

		result := collections.Map(input, func(i int) int {
			return i * 2
		})

		assert.Equal(t, expected, result)
	})

	t.Run("works with complex types", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		input := []Person{
			{Name: "Alice", Age: 30},
			{Name: "Bob", Age: 25},
		}
		expected := []string{"Alice", "Bob"}

		result := collections.Map(input, func(p Person) string {
			return p.Name
		})

		assert.Equal(t, expected, result)
	})
}
