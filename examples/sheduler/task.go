package main

import (
	"fmt"
	"time"

	"github.com/Hurricanezwf/toolbox/sheduler"
)

const (
	ShedulerTask = "task"
)

var (
	MaxConcurrency = 2
	QueueSize      = 4
)

func init() {
	if err := sheduler.Regist(&sheduler.Instance{
		Name: ShedulerTask,
		S:    sheduler.DefaultSheduler(),
		Conf: fmt.Sprintf(`{"maxconcurrency":%d,"queuesize":%d}`, MaxConcurrency, QueueSize),
	}); err != nil {
		msg := fmt.Sprintf("Regist task sheduler failed, %v", err)
		panic(msg)
	}
}

type Task struct {
	no int

	notify chan<- int
}

func NewTask(no int) *Task {
	return &Task{
		no: no,
	}
}

func (t *Task) SetNotify(notify chan<- int) *Task {
	t.notify = notify
	return t
}

func (t *Task) Do() {
	fmt.Printf("Hello, I'm No.%d\n", t.no)
	time.Sleep(time.Second)
	t.notify <- t.no
}
