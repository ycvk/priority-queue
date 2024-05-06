package priorityqueue

import (
	"math/rand"
	"testing"
)

func BenchmarkHeapPriorityQueue(b *testing.B) {
	pq := New[int, int](MinHeap)

	// 测试 Put
	b.Run("Put", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pq.Put(rand.Int(), rand.Int())
		}
	})

	// 填充优先级队列
	for i := 0; i < 1000; i++ {
		pq.Put(rand.Int(), rand.Int())
	}

	// 测试 Get
	b.Run("Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = pq.Get()
		}
	})

	// 测试 GetAndPop
	b.Run("GetAndPop", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = pq.GetAndPop()
			pq.Put(rand.Int(), rand.Int())
		}
	})

	// 测试 Update
	b.Run("Update", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			item := pq.GetAndPop()
			pq.Update(item.Value, rand.Int())
			pq.Put(item.Value, item.Priority)
		}
	})

	b.Run("upsert", func(b *testing.B) {
		for range b.N {
			pq.Upsert(rand.Int(), rand.Int())
		}
	})
}
