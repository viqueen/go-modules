package collections

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	t.Run("transforms slice of integers to strings", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		expected := []string{"1", "2", "3", "4", "5"}
		
		result := Map(input, func(i int) string {
			return strconv.Itoa(i)
		})
		
		assert.Equal(t, expected, result)
	})

	t.Run("transforms slice of strings to their lengths", func(t *testing.T) {
		input := []string{"hello", "world", "go", "test"}
		expected := []int{5, 5, 2, 4}
		
		result := Map(input, func(s string) int {
			return len(s)
		})
		
		assert.Equal(t, expected, result)
	})

	t.Run("handles empty slice", func(t *testing.T) {
		input := []int{}
		expected := []string{}
		
		result := Map(input, func(i int) string {
			return strconv.Itoa(i)
		})
		
		assert.Equal(t, expected, result)
	})

	t.Run("handles nil slice", func(t *testing.T) {
		var input []int
		
		result := Map(input, func(i int) string {
			return strconv.Itoa(i)
		})
		
		assert.Nil(t, result)
	})

	t.Run("preserves order", func(t *testing.T) {
		input := []int{3, 1, 4, 1, 5}
		expected := []int{6, 2, 8, 2, 10}
		
		result := Map(input, func(i int) int {
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
		
		result := Map(input, func(p Person) string {
			return p.Name
		})
		
		assert.Equal(t, expected, result)
	})
}