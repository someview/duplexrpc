package util

import (
	"fmt"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l1/tl.gobase.git/dt/pool"
	"sync/atomic"
)

// LinkedList 可并行读写的链表
type LinkedList[T any] struct {
	length int64

	head  *linkNode
	read  *linkNode
	write *linkNode
}

func NewLinkedList[T any]() *LinkedList[T] {
	l := &LinkedList[T]{}
	node := nodePool.Get()
	l.head = node
	l.read = node
	l.write = node
	return l
}

// Append 添加数据
func (l *LinkedList[T]) Append(value T) {
	node := l.write
	node.value = value
	l.growth()
	l.recalLen(1)
}

func (l *LinkedList[T]) growth() {
	if l.write.next == nil {
		l.write.next = nodePool.Get()
		l.write = l.write.next
	}
}

// Next 读取数据
func (l *LinkedList[T]) Next() (value T, err error) {
	if l.Len() <= 0 {
		return value, fmt.Errorf("empty")
	}
	l.recalLen(-1)

	read := l.read
	value, ok := read.value.(T)
	if !ok {
		fmt.Println("类型转换失败")
	}

	l.read = read.next

	return
}

// Len 链表长度
func (l *LinkedList[T]) Len() int {
	length := atomic.LoadInt64(&l.length)
	return int(length)
}

// Release 释放节点
func (l *LinkedList[T]) Release() {
	for l.head != l.read {
		node := l.head
		l.head = l.head.next
		node.recycle()
	}
}

func (l *LinkedList[T]) recalLen(delta int) (length int) {
	return int(atomic.AddInt64(&l.length, int64(delta)))
}

var nodePool = pool.NewSyncPool[*linkNode](func() any {
	return &linkNode{}
})

type linkNode struct {
	value any
	next  *linkNode
}

func (n *linkNode) recycle() {
	n.value = nil
	n.next = nil
	nodePool.Put(n)
}
