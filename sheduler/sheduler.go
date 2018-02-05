// sheduler实现了类似于协程池的概念，用于需要并发处理的任务进行并发量的控制
package sheduler

import (
	"bytes"
	"errors"
	"fmt"
	"time"
)

const (
	BlockForever time.Duration = 100 * 365 * 24 * time.Hour
)

var (
	ErrShedulerNotFound = errors.New("sheduler not found")
)

var (
	// 管理所有注册的sheduler实例
	// map[shedulerName]shedulerIns
	shedulers = make(map[string]*Instance)

	// sheduler包运行配置
	runConf *RunConf

	// 控制sheduler的监控
	monitorC chan struct{}
)

type RunConf struct {
	// 监控sheduler的时间间隔，如果小于等于0，将不启用监控
	MonitorInterval time.Duration

	// 向上层通知监控报告的channel，如果为空将不通知
	// 由sheduler.Run()的调用者负责释放MonitorReportC
	MonitorReportC chan<- string
}

type Instance struct {
	// sheduler名字
	Name string

	// sheduler的调度策略
	S Sheduler

	// sheduler配置，该字段会在Run的时候传进去，需要各个sheduler的实现进行处理
	// 推荐使用json字符串的形式传入
	Conf interface{}
}

// 默认调度策略
func DefaultSheduler() Sheduler {
	return &S1{}
}

func DefaultShedulerConf() interface{} {
	return DefaultS1Conf()
}

// 为不同的名字绑定不同的调度策略
func Regist(s *Instance) error {
	if s == nil {
		return errors.New("Nil shedulerIns")
	}
	if len(s.Name) <= 0 {
		return errors.New("Missing name")
	}
	if s.S == nil {
		return errors.New("Sheduler is nil")
	}
	if s.Conf == nil {
		return errors.New("Missing conf")
	}
	if _, ok := shedulers[s.Name]; ok {
		return errors.New("Sheduler has been existed!!!")
	}
	shedulers[s.Name] = s
	return nil
}

func Run(conf *RunConf) {
	if conf == nil {
		panic("Nil RunConf")
	}

	var err error
	for name, ins := range shedulers {
		if err = ins.S.Run(ins.Conf); err != nil {
			msg := fmt.Sprintf("Run %s sheduler failed, %v", name, err)
			panic(msg)
		}
	}

	go runMonitor(conf.MonitorInterval, conf.MonitorReportC)
}

func Close() {
	if monitorC != nil {
		close(monitorC)
		monitorC = nil
	}

	for _, ins := range shedulers {
		ins.S.Close()
	}
}

// Add add request to specified sheduler with timeout
func Add(shedulerName string, req Request, timeout time.Duration) error {
	ins, ok := shedulers[shedulerName]
	if !ok {
		return ErrShedulerNotFound
	}
	return ins.S.Add(req, timeout)
}

func runMonitor(interval time.Duration, report chan<- string) {
	if interval <= time.Duration(0) {
		return
	}

	if monitorC != nil {
		close(monitorC)
		time.Sleep(time.Second)
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("An exception occured: %v\n", err)
		}
	}()

	buf := bytes.NewBuffer(nil)
	monitorC = make(chan struct{})

	for {
		select {
		case <-monitorC:
			return
		case <-time.After(time.Duration(interval)):
			buf.Reset()
			buf.WriteString("\n----------- shedulers monitor report -----------\n")
			for name, ins := range shedulers {
				if ins == nil || ins.S == nil {
					continue
				}
				buf.WriteString(fmt.Sprintf("%-1s%-16s%-1s -- ", "[", name, "]"))
				buf.WriteString(ins.S.Monitor())
				buf.WriteString("\n")
			}
			buf.WriteString("------------------------------------------------\n")

			if report != nil {
				report <- buf.String()
			}
		}
	}
}
