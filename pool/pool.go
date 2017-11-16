package pool

type Concurrency struct {
	// max表示允许的最大并发量
	max int

	ch chan struct{}
}

func NewConcurrency(max int) *Concurrency {
	if max <= 0 {
		return nil
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

func (c *Concurrency) Get() {
	<-c.ch
}

func (c *Concurrency) Put() {
	c.ch <- struct{}{}
}
