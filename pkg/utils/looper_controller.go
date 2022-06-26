package utils

import "sync"

type LoopController struct {
	mu           sync.Mutex
	runningTasks int
}

func (l *LoopController) IncrementTask() {
	l.mu.Lock()
	l.runningTasks++
	l.mu.Unlock()
}

func (l *LoopController) DecrementTask() {
	l.mu.Lock()
	if l.runningTasks > 0 {
		l.runningTasks--
	}
	l.mu.Unlock()
}

func (l *LoopController) Running() bool {
	return l.runningTasks > 0
}
