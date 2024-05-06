package priorityqueue

import (
	"container/heap"
	"golang.org/x/exp/constraints"
	"sync/atomic"
	"unsafe"
)

type PriorityQueue[T any, P constraints.Ordered] interface {
	Len() int
	Put(value T, priority P)
	BatchPut(items ...*Item[T, P])
	Get() *Item[T, P]
	GetAndPop() *Item[T, P]
	IsEmpty() bool
	Update(value T, priority P)
	Clear()
}

// HeapType 指定堆类型 - 最小或最大
type HeapType int

const (
	MinHeap HeapType = iota // 小顶堆
	MaxHeap                 // 大顶堆
)

// Item 代表优先级队列中的一个元素
type Item[T any, P constraints.Ordered] struct {
	Value    T // 元素值
	Priority P // 元素优先级
}

// HeapPriorityQueue 基于容器/堆的优先级队列实现
type HeapPriorityQueue[T comparable, P constraints.Ordered] struct {
	items    unsafe.Pointer // *[]*Item[T, P]
	lookup   unsafe.Pointer // *map[T]int
	heapType HeapType
}

// New 创建一个新的优先级队列
func New[T comparable, P constraints.Ordered](kind HeapType) PriorityQueue[T, P] {
	items := make([]*Item[T, P], 0)
	lookup := make(map[T]int)
	pq := &HeapPriorityQueue[T, P]{
		items:    unsafe.Pointer(&items),
		lookup:   unsafe.Pointer(&lookup),
		heapType: kind,
	}
	return pq
}

// Len implements heap.Interface
func (pq *HeapPriorityQueue[T, P]) Len() int {
	items := atomic.LoadPointer(&pq.items)
	return len(*(*[]*Item[T, P])(items))
}

// Less implements heap.Interface
func (pq *HeapPriorityQueue[T, P]) Less(i, j int) bool {
	items := atomic.LoadPointer(&pq.items)
	slice := *(*[]*Item[T, P])(items)
	switch pq.heapType {
	case MinHeap:
		return slice[i].Priority < slice[j].Priority
	case MaxHeap:
		return slice[i].Priority > slice[j].Priority
	}
	return false // Should never reach here
}

// Swap implements heap.Interface
func (pq *HeapPriorityQueue[T, P]) Swap(i, j int) {
	items := atomic.LoadPointer(&pq.items)
	slice := *(*[]*Item[T, P])(items)
	slice[i], slice[j] = slice[j], slice[i]

	lookup := atomic.LoadPointer(&pq.lookup)
	lookupMap := *(*map[T]int)(lookup)
	lookupMap[slice[i].Value] = i
	lookupMap[slice[j].Value] = j
}

// Push implements heap.Interface
func (pq *HeapPriorityQueue[T, P]) Push(x any) {
	item := x.(*Item[T, P])
	items := (*[]*Item[T, P])(atomic.LoadPointer(&pq.items))
	n := len(*items)
	*items = append(*items, item)
	atomic.StorePointer(&pq.items, unsafe.Pointer(items))

	lookup := (*map[T]int)(atomic.LoadPointer(&pq.lookup))
	(*lookup)[item.Value] = n
	atomic.StorePointer(&pq.lookup, unsafe.Pointer(lookup))
}

// Pop implements heap.Interface 弹出最后一个元素
func (pq *HeapPriorityQueue[T, P]) Pop() any {
	items := (*[]*Item[T, P])(atomic.LoadPointer(&pq.items))
	n := len(*items)
	item := (*items)[n-1]
	*items = (*items)[:n-1]
	atomic.StorePointer(&pq.items, unsafe.Pointer(items))

	lookup := (*map[T]int)(atomic.LoadPointer(&pq.lookup))
	delete(*lookup, item.Value)
	atomic.StorePointer(&pq.lookup, unsafe.Pointer(lookup))

	return item
}

// Put 将元素添加到优先级队列中
func (pq *HeapPriorityQueue[T, P]) Put(value T, priority P) {
	item := &Item[T, P]{Value: value, Priority: priority}
	heap.Push(pq, item)
}

// Get 返回优先级队列中的下一个元素而不移除它
func (pq *HeapPriorityQueue[T, P]) Get() *Item[T, P] {
	items := atomic.LoadPointer(&pq.items)
	slice := *(*[]*Item[T, P])(items)
	if len(slice) == 0 {
		return nil
	}
	return slice[0]
}

// GetAndPop 移除并返回优先级队列中的下一个元素
func (pq *HeapPriorityQueue[T, P]) GetAndPop() *Item[T, P] {
	item := heap.Pop(pq)
	if item == nil {
		return nil
	}
	return item.(*Item[T, P])
}

// IsEmpty 检查优先级队列是否为空
func (pq *HeapPriorityQueue[T, P]) IsEmpty() bool {
	return pq.Len() == 0
}

// Update 更新元素的优先级
func (pq *HeapPriorityQueue[T, P]) Update(value T, priority P) {
	lookup := (*map[T]int)(atomic.LoadPointer(&pq.lookup))
	if index, ok := (*lookup)[value]; ok {
		items := (*[]*Item[T, P])(atomic.LoadPointer(&pq.items))
		if (*items)[index].Priority != priority {
			(*items)[index].Priority = priority
			heap.Fix(pq, index)
		}
	}
}

// Clear 清空优先级队列
func (pq *HeapPriorityQueue[T, P]) Clear() {
	items := make([]*Item[T, P], 0)
	lookup := make(map[T]int)
	atomic.StorePointer(&pq.items, unsafe.Pointer(&items))
	atomic.StorePointer(&pq.lookup, unsafe.Pointer(&lookup))
}

// BatchPut 批量将元素添加到优先级队列中
func (pq *HeapPriorityQueue[T, P]) BatchPut(items ...*Item[T, P]) {
	if len(items) == 0 {
		return
	}

	// 将新元素追加到切片末尾
	oldItems := (*[]*Item[T, P])(atomic.LoadPointer(&pq.items))
	newItems := append(*oldItems, items...)
	atomic.StorePointer(&pq.items, unsafe.Pointer(&newItems))

	// 更新lookup映射
	oldLookup := (*map[T]int)(atomic.LoadPointer(&pq.lookup))
	newLookup := make(map[T]int, len(*oldLookup)+len(items))
	for k, v := range *oldLookup {
		newLookup[k] = v
	}
	for i, item := range items {
		newLookup[item.Value] = len(*oldItems) + i
	}
	atomic.StorePointer(&pq.lookup, unsafe.Pointer(&newLookup))

	// 调整堆
	for i := len(*oldItems); i < len(newItems); i++ {
		heap.Fix(pq, i)
	}
}
