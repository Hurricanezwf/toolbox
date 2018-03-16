package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/Hurricanezwf/toolbox/pool"
)

func main() {
	runtime.GOMAXPROCS(4)
	p := pool.NewLimitPool(5).Fill()

	// put
	go func() {
		for {
			p.Put(struct{}{})
			fmt.Printf("[P] Producer\n")
			time.Sleep(time.Second)
		}
	}()

	// get
	go func() {
		for {
			p.Get()
			fmt.Printf("[C] Consumer\n")
			time.Sleep(time.Second)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	<-sig
}
