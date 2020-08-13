package crontab

import (
	"fmt"
	"time"

	"github.com/Hurricanezwf/toolbox/crontab/scheduler/types"
)

type TaskFunc func() error

type Task struct {
	// 任务名
	TaskName string

	// 任务到点后的执行函数
	DoFunc TaskFunc

	// crontab的执行时间配置
	specStr   string
	scheduler types.Scheduler

	// 计算任务下次执行时间的基准时间
	base time.Time
	// 任务下次执行时间
	next time.Time
}

// You should call SpecValid() to check if spec is valid before NewTask, or your program may be panic
func NewTask(
	tname string,
	specStr string,
	scheduler types.Scheduler,
	f TaskFunc,
) *Task {

	return &Task{
		TaskName:  tname,
		specStr:   specStr,
		scheduler: scheduler,
		DoFunc:    f,
		base:      time.Now().Local(),
	}
}

func (t *Task) ResetBaseTime(baseTime time.Time) *Task {
	t.base = baseTime
	return t
}

func (t *Task) doFuncCall() {
	t.DoFunc()
}

func (t *Task) setNext() {
	t.next = t.scheduler.Next(t.base)
	t.base = t.next
}

func (t *Task) String() string {
	return fmt.Sprintf("%-20s : %-20s", t.TaskName, t.specStr)
}
