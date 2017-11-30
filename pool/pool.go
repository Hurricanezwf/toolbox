package pool

import (
	"sync"
	"time"
)

// 保证并发数至少为1
type Concurrency struct {
	mutex  sync.Mutex
	curNum int // 当前并发数

	max int // max表示允许的最大并发量

	ch chan struct{}
}

func NewConcurrency(max int) *Concurrency {
	if max <= 0 {
		max = 1
	}

	ch := make(chan struct{}, max)
	for i := 0; i < max; i++ {
		ch <- struct{}{}
	}

	return &Concurrency{
		max: max,
		ch:  ch,
	}
}

// ResetMax: 重新设置最大并发量
func (c *Concurrency) ResetMax(max int) {
	if max <= 0 {
		max = 1
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.max == max {
		return
	}

	var (
		reput = 0
		newCh = make(chan struct{}, max)
	)

	if max > c.max {
		reput = max - c.curNum
		if reput < 0 {
			reput = 0
		}
		for reput > 0 {
			reput--
			newCh <- struct{}{}
		}
	} else if max < c.max {
		if c.curNum < max {
			reput = max - c.curNum
		}
		for reput > 0 {
			reput--
			newCh <- struct{}{}
		}
	}

	close(c.ch)
	c.ch = newCh
	c.max = max
}

func (c *Concurrency) Get() {
	for {
		if _, ok := <-c.ch; ok {
			break
		}
		time.Sleep(1 * time.Millisecond)
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.curNum++
}

func (c *Concurrency) Put() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.curNum--
	if c.curNum < 0 {
		c.curNum = len(c.ch)
	}
	if c.curNum < c.max {
		c.ch <- struct{}{}
	}
}
