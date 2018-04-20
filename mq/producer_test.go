package mq

import (
	"log"
	"testing"
	"time"
)

var (
	MQURL = "amqp://zwf:123456@localhost:5672/"
)

func TestProducer(t *testing.T) {
	m, err := New(MQURL).Open()
	if err != nil {
		t.Fatalf("%v", err)
	}

	p, err := m.Producer("test-producer")
	if err != nil {
		t.Fatalf("Create producer failed, %v", err)
	}

	exb := []*ExchangeBinds{
		&ExchangeBinds{
			Exch: DefaultExchange("exch.unitest", ExchangeDirect),
			Bindings: []*Binding{
				&Binding{
					RouteKey: "route.unitest1",
					Queues: []*Queue{
						DefaultQueue("queue.unitest1"),
					},
				},
				&Binding{
					RouteKey: "route.unitest2",
					Queues: []*Queue{
						DefaultQueue("queue.unitest2"),
					},
				},
			},
		},
	}

	if err = p.SetExchangeBinds(exb).Open(); err != nil {
		t.Fatalf("Open failed, %v", err)
	}

	for i := 0; i < 10; i++ {
		if i > 0 && i%3 == 0 {
			p.ch.Close() // 模拟channel关闭, 测试重联
		}
		err = p.Publish("exch.unitest", "route.unitest2", NewPublishMsg([]byte(`{"name":"zwf"}`)))
		log.Printf("Produce state:%d, err:%v\n", p.State(), err)
		time.Sleep(2 * time.Second)
	}

	p.Close()
	m.Close()
}
