package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/Hurricanezwf/toolbox/pool"
)

var (
	concurrency = 1
	c           = pool.NewConcurrency(concurrency)
)

func main() {
	go DoWork()

	go func() {
		http.HandleFunc("/", handler)
		http.ListenAndServe("127.0.0.1:9000", nil)
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGTERM)
	for {
		s := <-sig
		if s == syscall.SIGTERM {
			break
		}

		if s == syscall.SIGUSR1 {
			fmt.Printf("Receive add concurrency\n")
			concurrency += 2
			c.ResetMax(concurrency)
		} else if s == syscall.SIGUSR2 {
			fmt.Printf("Receive reduce concurrency\n")
			concurrency -= 3
			c.ResetMax(concurrency)
		} else {
			fmt.Printf("Unkonw signal %v\n", s)
		}
	}
}

func DoWork() {
	for {
		c.Get()
		go Do(c)
	}
}

func Do(c *pool.Concurrency) {
	defer func() {
		c.Put()
	}()

	time.Sleep(5 * time.Second)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	p := pprof.Lookup("goroutine")
	p.WriteTo(w, 1)
}
