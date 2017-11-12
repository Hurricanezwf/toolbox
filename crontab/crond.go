package crontab

import (
	"errors"
	"sync"
	"time"
	"toolbox/crontab/heap"
)

var (
	forever = 100 * 365 * 24 * time.Hour
)

type Crond struct {
	mutex sync.RWMutex
	h     *heap.MinHeap

	changeC chan struct{}

	running bool
	stopC   chan struct{}
}

func (c *Crond) Run() {
	if c.running {
		return
	}

	c.running = true

	c.h = heap.NewMinHeap(1024)
	c.changeC = make(chan struct{}, 1)
	c.stopC = make(chan struct{})

	for {
		sleep := func() time.Duration {
			waitTime := forever
			start := time.Now()
			for {
				e := c.GetTop()
				if e == nil {
					break
				}

				w, expired := e.Sub(time.Now())
				if !expired {
					waitTime = w
					break
				}

				// for all expired tasks in the top, calc next exec time and repush into heap
				c.PopAndRepush(false)
			}
			fix := time.Since(start)
			return waitTime - fix
		}()

		select {
		case <-c.stopC:
			return
		case <-c.changeC:
			continue
		case <-time.After(sleep):
			c.PopAndRepush(true)
			continue
		}
	}
}

func (c *Crond) Close() {
	close(c.stopC)
	c.running = false
}

func (c *Crond) Add(t *Task) error {
	if t == nil {
		return errors.New("nil task")
	}

	// 计算该任务最近一次执行时间
	if t.next.IsZero() {
		t.setNext()
	}

	// 1. 先尝试merge
	var setCount int
	c.mutex.Lock()
	c.h.Walk(func(e heap.Element) heap.Element {
		te, ok := e.(*TaskElement)
		if !ok {
			return nil
		}
		if !t.next.Equal(te.next) {
			return nil
		}
		te.SetTask(t)
		setCount++
		return te
	})
	c.mutex.Unlock()
	if setCount > 0 {
		c.changeC <- struct{}{}
		return nil
	}

	// 2. heap中无同一时刻执行的任务元素，新建一个
	taskElem := NewTaskElement()
	taskElem.SetNext(t.next)
	taskElem.SetTask(t)
	c.mutex.Lock()
	err := c.h.Push(taskElem)
	c.mutex.Unlock()
	if err != nil {
		return err
	}

	c.changeC <- struct{}{}

	return nil
}

func (c *Crond) Del(taskName string) int {
	var affected int
	c.mutex.Lock()
	c.h.Walk(func(e heap.Element) heap.Element {
		te, ok := e.(*TaskElement)
		if !ok {
			return nil
		}
		if _, exist := te.tasks[taskName]; !exist {
			return nil
		}
		delete(te.tasks, taskName)
		affected++
		return te
	})
	c.mutex.Unlock()
	if affected > 0 {
		c.changeC <- struct{}{}
	}
	return affected
}

func (c *Crond) GetTop() *TaskElement {
	c.mutex.RLock()
	e := c.h.Top()
	c.mutex.RUnlock()
	if e == nil {
		return nil
	}
	taskElem, ok := e.(*TaskElement)
	if !ok {
		return nil
	}
	return taskElem
}

func (c *Crond) PopAndRepush(doFunc bool) {
	c.mutex.Lock()
	e := c.h.Pop().(*TaskElement)
	c.mutex.Unlock()
	for _, t := range e.Value().(map[string]*Task) {
		if doFunc {
			go t.doFuncCall()
		}
		t.setNext()
		c.Add(t)
	}
}
