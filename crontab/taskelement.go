package crontab

import (
	"time"
)

// TaskElement implements of heap.Element
type TaskElement struct {
	next  time.Time
	tasks map[string]*Task
}

func NewTaskElement() *TaskElement {
	return &TaskElement{
		tasks: make(map[string]*Task),
	}
}

func (e *TaskElement) SetNext(next time.Time) {
	e.next = next
}

func (e *TaskElement) SetTask(t *Task) {
	e.tasks[t.TaskName] = t
}

func (e *TaskElement) Key() interface{} {
	return e.next
}

func (e *TaskElement) Value() interface{} {
	return e.tasks
}

func (e *TaskElement) Compare(next interface{}) int {
	if e.next.Equal(next.(time.Time)) {
		return 0
	}
	if e.next.Before(next.(time.Time)) {
		return -1
	}
	return 1
}

func (e *TaskElement) Sub(t time.Time) (du time.Duration, expired bool) {
	if e.next.Before(t) {
		return du, true
	}
	return e.next.Sub(t), false
}
