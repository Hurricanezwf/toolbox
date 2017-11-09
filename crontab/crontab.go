// 迅速添加任务
// 不太慢的删除任务
// 添加任务对其他任务执行干扰小
package crontab

import "time"

func Add(t *Task) {

}

//////////////////////////////////////////////////
type TaskFunc func(param interface{}) error

type Task struct {
	// 任务名
	TaskName string

	// 任务到点后的执行函数
	DoFunc    TaskFunc
	FuncParam interface{}

	// crontab的执行时间配置
	Spec string

	// 任务的开始时间,可以是创建的时间或者暂停后重新开始的时间
	StartTime int64

	// 任务最近一次执行时间
	LastExecTime int64
}

func NewTask(tname, spec string, f TaskFunc, fParam interface{}) *Task {
	return &Task{
		TaskName:     tname,
		Spec:         spec,
		DoFunc:       f,
		FuncParam:    fParam,
		StartTime:    time.Now().Unix(),
		LastExecTime: 0,
	}
}
