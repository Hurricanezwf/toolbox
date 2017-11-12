package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
	"toolbox/crontab"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	crontab.Open()

	t1 := crontab.NewTask("test1", "0/1 * * * * *", Do1, nil)
	crontab.Add(t1)

	t2 := crontab.NewTask("test2", "0/5 * * * * *", Do2, nil)
	crontab.Add(t2)

	time.Sleep(5 * time.Second)
	aff := crontab.Del("test1")
	log.Printf("Affeted num: %d\n", aff)

	time.Sleep(5 * time.Second)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	crontab.Close()
}

func Do1(param interface{}) error {
	log.Println("Haha")
	return nil
}

func Do2(param interface{}) error {
	log.Println("Hello")
	return nil
}
