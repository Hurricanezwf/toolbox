package main

import (
	"fmt"
	"time"

	"github.com/Hurricanezwf/toolbox/sheduler"
)

func main() {
	// run sheduler
	sheduler.Run(&sheduler.RunConf{
		MonitorInterval: 0,
		MonitorReportC:  nil,
	})
	defer sheduler.Close()

	count := 10
	notify := make(chan int, 1)
	defer close(notify)

	start := time.Now()

	// produce tasks
	go func() {
		for i := 0; i < count; i++ {
			t := NewTask(i).SetNotify(notify)
			err := sheduler.Add(ShedulerTask, t, sheduler.BlockForever)
			if err != nil {
				panic(err)
			}
		}
	}()

	// consume task notify
	for i := 0; i < count; i++ {
		fmt.Printf("No.%d Done\n", <-notify)
	}

	fmt.Printf("All done, elapse:%v", time.Since(start))
}
