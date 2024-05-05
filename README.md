# 使用小/大顶堆实现的优先级队列

这是一个使用Go语言实现的通用优先级队列库，支持小顶堆和大顶堆两种模式。该库基于Go的`container/heap`包，提供了一种高效的方式来管理和操作具有优先级的元素。

## 特性

- 支持泛型，允许使用任意可比较类型作为元素值和优先级
- 提供小顶堆和大顶堆两种模式
- 线程安全，使用读写锁保护共享资源
- 高效的堆操作，时间复杂度为O(log n)
- 支持更新元素的优先级
- 提供了清空优先级队列的方法

## 安装

使用Go模块管理依赖，添加以下内容到你的`go.mod`文件：

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
item := pq.Get() // 获取优先级最高的元素，但不移除
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
    mpq := pq.New[string, int](priorityqueue.MinHeap)

    mpq.Put("task1", 5)
    mpq.Put("task2", 3)
    mpq.Put("task3", 7)

    fmt.Println(mpq.GetAndPop()) // 输出：&{task2 3}
    fmt.Println(mpq.GetAndPop()) // 输出：&{task1 5}

    mpq.Update("task3", 1)

    fmt.Println(mpq.GetAndPop()) // 输出：&{task3 1}
}
```

## 许可证

本项目采用MIT许可证，详情请参阅[LICENSE](https://github.com/ycvk/priority-queue/blob/main/LICENSE)文件。

## 贡献

欢迎提交问题和合并请求。如果你发现任何bug或有任何改进建议，请随时提出。