package types

import "time"

// ParserType 解析器类型
type ParserType int

const (
	ParserBasic = iota // basic只具备了定制的定时任务表达式解析能力
	ParserPro          // pro具备完整的定时任务表达式解析能力
)

// Scheduler 定时任务调度器的抽象
type Scheduler interface {
	// Prev 找寻基准时间之前理论上的执行时间
	// pro类型的 Scheduler 没有实现 Prev, basic类型的 Scheduler 实现了
	Prev(baseTime time.Time) time.Time

	// Next 找寻基准时间之后理论上的执行时间
	Next(baseTime time.Time) time.Time
}

// Parser 表达式解析器的抽象
type Parser interface {
	Parse(spec string) (Scheduler, error)
}
