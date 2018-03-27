package sheduler

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"
)

type S1Conf struct {
	// 最大并发量, 默认为runtime.NumCPU()
	MaxConcurrency int `json:"maxconcurrency"`

	// 队列大小, 默认1024
	QueueSize int `json:"queuesize"`
}

func DefaultS1Conf() *S1Conf {
	return &S1Conf{
		MaxConcurrency: runtime.NumCPU(),
		QueueSize:      1024,
	}
}

func (c *S1Conf) Validate() error {
	if c.MaxConcurrency <= 0 {
		return errors.New("S1Conf: Missing 'MaxConcurrency'")
	}
	if c.MaxConcurrency > 10000 {
		return errors.New("S1Conf: Too large 'Maxconcurrency', max is 10000")
	}

	if c.QueueSize < 0 {
		return errors.New("S1Conf: Missing 'QueueSize'")
	}
	if c.QueueSize > 50000 {
		return errors.New("S1Conf: Too large 'QueueSize', max is 50000")
	}
	return nil
}

////////////////////////////////////////////////////////////
type S1 struct {
	queueMutex sync.Mutex
	queue      chan Request

	conf *S1Conf

	stopC chan struct{}
}

func (s *S1) Run(conf interface{}) error {
	if conf == nil {
		return errors.New("Missing conf")
	}

	var (
		err     error
		tmpConf *S1Conf
	)

	switch t := conf.(type) {
	case []byte:
		tmpConf = &S1Conf{}
		if err = json.Unmarshal(conf.([]byte), tmpConf); err != nil {
			return fmt.Errorf("Bad conf format, %v", err)
		}
	case string:
		confBytes := []byte(conf.(string))
		tmpConf = &S1Conf{}
		if err := json.Unmarshal(confBytes, tmpConf); err != nil {
			return fmt.Errorf("Bad conf format, %v", err)
		}
	case *S1Conf:
		tmpConf = conf.(*S1Conf)
	case S1Conf:
		var cfg = conf.(S1Conf)
		tmpConf = &cfg
	default:
		return fmt.Errorf("Invalid conf type %v", t)
	}

	if err = tmpConf.Validate(); err != nil {
		return err
	}

	s.conf = tmpConf
	s.queue = make(chan Request, s.conf.QueueSize)
	s.stopC = make(chan struct{})

	for i := 0; i < s.conf.MaxConcurrency; i++ {
		go s.run()
	}
	return nil
}

// TODO: 关闭时停止Add操作并等待队列为空
func (s *S1) Close() {
	close(s.stopC)
	s.stopC = nil

	close(s.queue)
	s.queue = nil
}

func (s *S1) run() {
	for {
		select {
		case <-s.stopC:
			return
		case req := <-s.queue:
			if req != nil {
				req.Do()
			}
		}
	}
}

func (s *S1) Add(req Request, timeout time.Duration) error {
	select {
	case s.queue <- req:
		return nil
	case <-time.After(timeout):
		return errors.New("timeout")
	}
	return nil
}

func (s S1) Monitor() string {
	return fmt.Sprintf("queue len: %d, cap: %d", len(s.queue), cap(s.queue))
}
