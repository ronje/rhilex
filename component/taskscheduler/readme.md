<!--
 Copyright (C) 2024 wwhai

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU Affero General Public License as
 published by the Free Software Foundation, either version 3 of the
 License, or (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU Affero General Public License for more details.

 You should have received a copy of the GNU Affero General Public License
 along with this program.  If not, see <https://www.gnu.org/licenses/>.
-->

# 轻量级定时调度器

## 概述

轻量级定时调度器是一种用于定时执行任务的工具，支持一次性任务和永久任务的调度。调度器采用任务槽的设计模式，通过定期检查任务槽，执行到期的任务。它允许用户传入上下文参数，以便对任务的执行进行控制和取消。

## 设计思路

### 1. 任务管理

调度器使用一个 `Task` 结构体来表示每个任务。每个任务包含以下信息：
- **ID**: 唯一标识符。
- **Duration**: 任务的执行间隔时间。
- **Callback**: 任务执行时调用的回调函数。
- **IsPermanent**: 标记任务是否为永久任务。
- **NextRun**: 下次运行时间。
- **Context**: 任务的上下文，支持取消和超时。
- **CancelFunc**: 用于取消上下文的函数。

### 2. 任务调度

调度器使用一个 `Scheduler` 结构体来管理任务。该结构体包含：
- **tasks**: 存储当前所有任务的映射。
- **taskID**: 用于生成唯一任务 ID 的计数器。
- **mu**: 用于保证并发安全的互斥锁。
- **stopChan**: 用于停止调度器的信号通道。
- **ticker**: 定时器，定期检查任务槽。

### 3. 定时执行

调度器内部有一个 goroutine，不断地从任务槽中检查任务的到期情况。当任务到期时，执行其回调函数，并根据任务的类型决定是否重新安排任务的下次执行时间或删除任务。

### 4. 上下文支持

每个任务的回调函数都接收一个 `context.Context` 参数，使得用户可以控制任务的执行，例如取消任务。调度器在删除任务时会调用取消函数，以确保上下文的资源得到释放。

## 使用说明

### 1. 创建调度器

要使用调度器，首先需要创建一个新的调度器实例：

```go
scheduler := NewScheduler()
```

### 2. 添加任务

可以通过 `AddTask` 方法添加新的任务。该方法接受三个参数：倒计时时间、回调函数和任务类型（一次性或永久）：

```go
taskID := scheduler.AddTask(2*time.Second, sampleTask, false) // 一次性任务
permanentTaskID := scheduler.AddTask(1*time.Second, sampleTask, true) // 永久任务
```

### 3. 删除任务

使用 `DeleteTask` 方法可以根据任务 ID 删除任务。删除任务时会自动取消其上下文：

```go
scheduler.DeleteTask(permanentTaskID)
```

### 4. 停止调度器

在不再需要调度器时，可以调用 `Stop` 方法来停止调度器并释放相关资源：

```go
scheduler.Stop()
```

## 示例代码

以下是一个完整的示例，展示如何使用轻量级定时调度器：

```go
func sampleTask(ctx context.Context) {
	select {
	case <-ctx.Done():
		fmt.Println("Task cancelled")
		return
	default:
		fmt.Println("Task executed at:", time.Now())
	}
}

func main() {
	scheduler := NewScheduler()
	scheduler.AddTask(2*time.Second, sampleTask, false)
	taskID := scheduler.AddTask(1*time.Second, sampleTask, true)
	time.Sleep(5 * time.Second)
	scheduler.DeleteTask(taskID)
	time.Sleep(3 * time.Second)
	scheduler.Stop()
}
```
