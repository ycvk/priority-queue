# 使用小/大顶堆实现的无锁化优先级队列

这是一个使用Go语言实现的通用优先级队列库,支持小顶堆和大顶堆两种模式。该库基于Go的`container/heap`包,提供了一种高效的方式来管理和操作具有优先级的元素。

## 特性

- 支持泛型,允许使用任意可比较类型作为元素值和优先级
- 提供小顶堆和大顶堆两种模式
- 无锁化设计,使用CAS原子操作实现并发安全,提高了并发性能
- 高效的堆操作,时间复杂度为O(log n)
- 支持更新元素的优先级
- 提供了清空优先级队列的方法

## 安装

使用Go模块管理依赖,添加以下内容到你的`go.mod`文件：

```go
import "github.com/ycvk/priority-queue"
```

然后运行`go mod tidy`下载依赖。

## 使用方法

### 创建优先级队列

```go
pq := New[int, int](MinHeap) // 创建一个小顶堆
pq := New[string, float64](MaxHeap) // 创建一个大顶堆
```

### 添加元素

```go
pq.Put(10, 5)
pq.Put(20, 3)
pq.Put(30, 7)
```

### 获取元素

```go
item := pq.Get() // 获取优先级最高的元素,但不移除
item = pq.GetAndPop() // 获取并移除优先级最高的元素
```

### 更新元素的优先级

```go
pq.Update(20, 8) // 更新元素20的优先级为8
```

### 清空优先级队列

```go
pq.Clear()
```

## 示例

```go
package main

import (
    "fmt"
    pq "github.com/ycvk/priority-queue"
)

func main() {
    mpq := pq.New[string, int](pq.MinHeap)

    mpq.Put("task1", 5)
    mpq.Put("task2", 3)
    mpq.Put("task3", 7)

    fmt.Println(mpq.GetAndPop()) // 输出：&{task2 3}
    fmt.Println(mpq.GetAndPop()) // 输出：&{task1 5}

    mpq.Update("task3", 1)

    fmt.Println(mpq.GetAndPop()) // 输出：&{task3 1}
}
```

## 无锁化设计

本库采用无锁化设计,使用CAS(Compare-And-Swap)原子操作来实现并发安全。相比传统的互斥锁,无锁化设计可以显著提高并发性能,减少线程阻塞和上下文切换的开销。

在无锁化实现中,我使用`unsafe.Pointer`和原子操作来管理共享资源,如切片和映射。通过原子地加载和存储指针,保证了并发访问的安全性。

同时,我优化了一些关键操作,如`Push`和`Pop`,通过直接操作切片指针来避免不必要的内存分配和复制。在`Update`操作中,只有当元素的优先级实际发生变化时,才触发堆调整,减少了冗余操作。

这些优化和无锁化设计使得本库在高并发场景下表现出色,同时保证了数据的正确性和一致性。

## Benchmark
使用`mutex`加锁方式时(代码详见`v1_mutex`目录下):
```
goos: darwin
goarch: arm64
pkg: github.com/ycvk/priority-queue
BenchmarkHeapPriorityQueue
BenchmarkHeapPriorityQueue/Put
BenchmarkHeapPriorityQueue/Put-10         	 2501148	       589.9 ns/op
BenchmarkHeapPriorityQueue/Get
BenchmarkHeapPriorityQueue/Get-10         	48504936	        24.67 ns/op
BenchmarkHeapPriorityQueue/GetAndPop
BenchmarkHeapPriorityQueue/GetAndPop-10   	  253447	      4601 ns/op
BenchmarkHeapPriorityQueue/Update
BenchmarkHeapPriorityQueue/Update-10      	  423760	      2821 ns/op
PASS
Exiting.
```
使用目前CAS原子化操作时:
```
goos: darwin
goarch: arm64
pkg: github.com/ycvk/priority-queue
BenchmarkHeapPriorityQueue
BenchmarkHeapPriorityQueue/Put
BenchmarkHeapPriorityQueue/Put-10         	 4053801	       369.4 ns/op
BenchmarkHeapPriorityQueue/Get
BenchmarkHeapPriorityQueue/Get-10         	582441570	         2.057 ns/op
BenchmarkHeapPriorityQueue/GetAndPop
BenchmarkHeapPriorityQueue/GetAndPop-10   	  354594	      2124 ns/op
BenchmarkHeapPriorityQueue/Update
BenchmarkHeapPriorityQueue/Update-10      	 1000000	      1045 ns/op
PASS
```

可以看到快了几倍以上。
## 许可证

本项目采用许可证详情请参阅[LICENSE](https://github.com/ycvk/priority-queue/blob/main/LICENSE)文件。

## 贡献

欢迎提交问题和合并请求。如果你发现任何bug或有任何改进建议,请随时提出。我们非常重视性能和并发安全,如果你有任何优化思路或发现潜在的并发问题,也欢迎与我们分享。
