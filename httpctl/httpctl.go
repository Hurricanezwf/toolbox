package httpctl

import (
	"errors"
	"time"
)

var (
	ErrTimeout = errors.New("HTTPCtl: Block timeout")
	ErrClosed  = errors.New("HTTPCtl: Closed")
)

type HTTPCtl struct {
	// HTTP请求并发量限制
	limit int

	// 阻塞超时
	blockTimeout time.Duration

	// 资源池
	pool chan struct{}
}

func New(interval time.Duration, limit int) *HTTPCtl {
	if limit <= 0 {
		panic("Invalid limit value")
	}

	ctl := &HTTPCtl{
		limit:        limit,
		blockTimeout: 30 * time.Second,
		pool:         make(chan struct{}, limit),
	}

	// 定时补充资源
	// 截取当前时间点时剩余的资源从而计算出需要补充的资源，然后进行填充
	go func() {
		for {
			remain := len(ctl.pool)
			for add := ctl.limit - remain; add > 0; add-- {
				ctl.pool <- struct{}{}
			}
			time.Sleep(interval)
		}
	}()

	return ctl
}

func (ctl *HTTPCtl) SetBlockTimeout(timeout time.Duration) *HTTPCtl {
	ctl.blockTimeout = timeout
	return ctl
}

type RequestFunc func()

func (ctl *HTTPCtl) Do(f RequestFunc) error {
	select {
	case _, ok := <-ctl.pool:
		if !ok {
			return ErrClosed
		}
		// 并发调用
		go f()
	case <-time.After(ctl.blockTimeout):
		return ErrTimeout
	}
	return nil
}
