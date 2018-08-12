package main

import (
	"math/rand"
	"time"

	"github.com/Hurricanezwf/toolbox/httpctl"
)

func main() {
	ctl := httpctl.New(time.Second, 3)

	count := 1000
	rp := make(chan struct{}, 20)
	rand.Seed(time.Now().UnixNano())

	go func() {
		for i := 0; i < count; i++ {
			err := ctl.Do(func() {
				time.Sleep(time.Duration(rand.Intn(5000)) * time.Millisecond)
				rp <- struct{}{}
			})
			if err != nil {
				panic(err.Error())
			}
		}
	}()

	for ; count > 0; count-- {
		<-rp
	}
}
