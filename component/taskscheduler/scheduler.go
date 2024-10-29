// Copyright (C) 2024 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package taskscheduler

import (
	"context"
	"sync"
	"time"
)

type Task struct {
	ID          int
	Duration    time.Duration
	Callback    func(ctx context.Context)
	IsPermanent bool
	NextRun     time.Time
	Context     context.Context
	CancelFunc  context.CancelFunc
}

type Scheduler struct {
	tasks    map[int]*Task
	taskID   int
	mu       sync.Mutex
	stopChan chan bool
	ticker   *time.Ticker
}

// NewScheduler 创建新的调度器
func NewScheduler() *Scheduler {
	s := &Scheduler{
		tasks:    make(map[int]*Task),
		stopChan: make(chan bool),
		ticker:   time.NewTicker(100 * time.Millisecond), // 每100毫秒检查一次
	}
	go s.run()
	return s
}

// AddTask 添加一个新的任务
func (s *Scheduler) AddTask(duration time.Duration, callback func(ctx context.Context), isPermanent bool) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.taskID++
	ctx, cancel := context.WithCancel(context.Background())

	task := &Task{
		ID:          s.taskID,
		Duration:    duration,
		Callback:    callback,
		IsPermanent: isPermanent,
		NextRun:     time.Now().Add(duration),
		Context:     ctx,
		CancelFunc:  cancel,
	}

	s.tasks[task.ID] = task
	return task.ID
}

// DeleteTask 删除一个任务
func (s *Scheduler) DeleteTask(taskID int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task, exists := s.tasks[taskID]; exists {
		task.CancelFunc()
		delete(s.tasks, taskID)
	}
}

// run 监视任务并执行到期的任务
func (s *Scheduler) run() {
	for {
		select {
		case <-s.stopChan:
			s.ticker.Stop()
			return
		case <-s.ticker.C:
			s.executeDueTasks()
		}
	}
}

// executeDueTasks 执行到期的任务
func (s *Scheduler) executeDueTasks() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for id, task := range s.tasks {
		if task.NextRun.Before(now) || task.NextRun.Equal(now) {
			task.Callback(task.Context)
			if task.IsPermanent {
				task.NextRun = now.Add(task.Duration)
			} else {
				task.CancelFunc()
				delete(s.tasks, id)
			}
		}
	}
}

// 停止调度器
func (s *Scheduler) Stop() {
	s.stopChan <- true
}
