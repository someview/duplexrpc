package util

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestLinkedList_Append_Next(t *testing.T) {
	link := NewLinkedList[int]()
	num := 3
	for i := 0; i < num; i++ {
		link.Append(i)
	}

	for i := 0; i < num; i++ {
		n, _ := link.Next()
		assert.Equal(t, i, n)
	}

	assert.Equal(t, 0, link.Len())

}

func TestLinkedList_Release(t *testing.T) {
	link := NewLinkedList[int]()
	num := 100
	for i := 0; i < num; i++ {
		link.Append(i)
	}
	link.Release()
	assert.Equal(t, num, link.Len())
	link.Next()
	link.Release()
	assert.Equal(t, num-1, link.Len())
	for link.Len() > 0 {
		link.Next()
	}
	assert.Equal(t, 0, link.Len())
	link.Release()
	assert.Equal(t, link.head, link.read)
	assert.Equal(t, link.head, link.write)

}

func TestLinkedList_RWRace(t *testing.T) {
	// use go test -race
	num := 100000
	wg := sync.WaitGroup{}
	wg.Add(num)
	link := NewLinkedList[int]()
	go func() {
		for i := 0; i < num; i++ {
			link.Append(i)
		}
	}()
	go func() {
		for i := 0; i < num; i++ {
			for {
				value, err := link.Next()
				if err != nil {
					continue
				}
				assert.Equal(t, i, value)
				break
			}
			wg.Done()
		}
	}()
	wg.Wait()
}
