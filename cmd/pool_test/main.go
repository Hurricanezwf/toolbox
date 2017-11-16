package main

import (
	"fmt"
	"time"

	"github.com/Hurricanezwf/toolbox/pool"
)

func main() {
	c := pool.NewConcurrency(2)

	for i := 0; i < 100; i++ {
		c.Get()
		go Do(i, c)
	}
}

func Do(idx int, c *pool.Concurrency) {
	fmt.Printf("%d is running...\n", idx)
	defer func() {
		fmt.Printf("%d stopped\n", idx)
		c.Put()
	}()

	time.Sleep(3 * time.Second)
}
