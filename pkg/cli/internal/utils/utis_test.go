package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdvancedArray_Next(t *testing.T) {
	t.Parallel()

	arr := []int{1, 2, 3}
	advArray := NewAdvancedArray(arr)

	val, ok := advArray.Next()
	assert.True(t, ok)
	assert.Equal(t, 1, val)

	val, ok = advArray.Next()
	assert.True(t, ok)
	assert.Equal(t, 2, val)

	val, ok = advArray.Next()
	assert.True(t, ok)
	assert.Equal(t, 3, val)

	val, ok = advArray.Next()
	assert.False(t, ok)
	assert.Equal(t, 0, val)
}

func TestAdvancedArray_EmptyArray(t *testing.T) {
	t.Parallel()

	arr := []int{}
	advArray := NewAdvancedArray(arr)

	val, ok := advArray.Next()
	assert.False(t, ok)
	assert.Equal(t, 0, val)
}

func TestAdvancedArray_Back(t *testing.T) {
	t.Parallel()

	arr := []int{1, 2, 3}
	advArray := NewAdvancedArray(arr)

	// Move forward in the array
	advArray.Next()
	advArray.Next()

	// Move back and check the value
	advArray.Back()
	val, ok := advArray.Next()
	assert.True(t, ok)
	assert.Equal(t, 2, val)

	// Move back to the start and check the value
	advArray.Back()
	val, ok = advArray.Next()
	assert.True(t, ok)
	assert.Equal(t, 2, val)

	// Attempt to move back beyond the start
	advArray.Back()
	advArray.Back()
	val, ok = advArray.Next()
	assert.True(t, ok)
	assert.Equal(t, 1, val)
}

func TestAdvancedArray_BackOnEmptyArray(t *testing.T) {
	t.Parallel()

	arr := []int{}
	advArray := NewAdvancedArray(arr)

	// Attempt to move back on an empty array
	advArray.Back()
	val, ok := advArray.Next()
	assert.False(t, ok)
	assert.Equal(t, 0, val)
}
