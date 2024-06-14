package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderedMap_FrontAndBack(t *testing.T) {
	om := NewOrderedMap[int, string]()
	assert.Nil(t, om.Front())
	assert.Nil(t, om.Back())
	om.Set(1, "1")
	om.Set(2, "2")
	assert.Equal(t, om.Front().Value, "1")
	assert.Equal(t, om.Back().Value, "2")
}

func TestOrderedMap_Next(t *testing.T) {
	om := NewOrderedMap[int, string]()
	om.Set(1, "1")
	om.Set(2, "2")
	v1, ok := om.Next(1)
	assert.True(t, ok, "存在")
	assert.Equal(t, v1, "2")
	v2, ok := om.Next(2)
	assert.True(t, ok, "存在")
	assert.Equal(t, v2, "1")
}

func TestOrderedMap_Keys(t *testing.T) {
	om := NewOrderedMap[int, string]()
	om.Set(1, "1")
	om.Set(2, "2")
	keys := om.Keys()
	assert.Equal(t, 2, len(keys))
	assert.Equal(t, keys[0], 1)
	assert.Equal(t, keys[1], 2)
}

func TestOrderedMap_Range(t *testing.T) {
	om := NewOrderedMap[int, int]()
	om.Set(1, 1)
	om.Set(2, 2)
	om.Set(3, 3)
	var count = 0
	om.Range(func(k, v int) bool {
		count++
		assert.Equal(t, k, v)
		return true
	})
	assert.Equal(t, 3, count)
}
