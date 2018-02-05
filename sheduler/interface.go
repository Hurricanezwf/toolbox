package sheduler

import "time"

// 调度器请求接口, 与models层解耦
type Request interface {
	Do()
}

type Sheduler interface {
	Run(conf interface{}) error
	Monitor() string
	Add(req Request, timeout time.Duration) error
	Close()
}
