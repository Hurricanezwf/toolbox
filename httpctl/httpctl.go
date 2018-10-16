package httpctl

import (
	"errors"
	"math"
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

	// 缓冲队列，用于存放受控制的名额
	cachedQueue chan struct{}

	// 正式队列，用于存放实际可消费的名额
	queue chan struct{}

	// 控制HTTPCtl关闭
	stopC chan struct{}
}

func New(interval time.Duration, limit int) *HTTPCtl {
	if limit <= 0 {
		panic("Invalid limit value")
	}

	ctl := &HTTPCtl{
		limit:        limit,
		blockTimeout: time.Duration(math.MaxInt64),
		cachedQueue:  make(chan struct{}, limit),
		queue:        make(chan struct{}, limit),
		stopC:        make(chan struct{}),
	}

	// 填充资源额度
	for ; limit > 0; limit-- {
		ctl.queue <- struct{}{}
	}

	// 定时将请求额度从缓冲队列搬到正式队列
	go func() {
		for {
			ticker := time.NewTicker(interval)
			defer ticker.Stop()

			select {
			case <-ctl.stopC:
				return
			case <-ticker.C:
				for cnt := len(ctl.cachedQueue); cnt > 0; cnt-- {
					ctl.queue <- <-ctl.cachedQueue
				}
			}
		}
	}()

	return ctl
}

func (ctl *HTTPCtl) Close() {
	close(ctl.stopC)
}

func (ctl *HTTPCtl) SetBlockTimeout(timeout time.Duration) *HTTPCtl {
	ctl.blockTimeout = timeout
	return ctl
}

type RequestFunc func()

func (ctl *HTTPCtl) Do(f RequestFunc) error {
	select {
	case <-ctl.stopC:
		return ErrClosed
	case v := <-ctl.queue:
		// 并发调用
		go func() {
			//fmt.Printf("[%d]\n", time.Now().Unix())
			f()
			ctl.cachedQueue <- v // 使用完成后统一丢入缓冲队列，后续统一刷入正式队列
		}()
	case <-time.After(ctl.blockTimeout):
		return ErrTimeout
	}
	return nil
}
