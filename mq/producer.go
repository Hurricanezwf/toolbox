package mq

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Hurricanezwf/toolbox/log"

	"github.com/streadway/amqp"
)

type Producer struct {
	// Name to identify the producer, "" is OK
	name string

	// MQ实例
	mq *MQ

	// 保护数据并发安全
	mutex sync.RWMutex

	// MQ的会话channel
	ch *amqp.Channel

	// MQ的exchange与其绑定的queues
	exchangeBinds []*ExchangeBinds

	// 监听会话channel关闭
	closeC chan *amqp.Error
	// Producer关闭控制
	stopC chan struct{}

	// Producer 状态
	state uint8
}

func newProducer(name string, mq *MQ) *Producer {
	return &Producer{
		name:  name,
		mq:    mq,
		stopC: make(chan struct{}),
		state: StateClosed,
	}
}

func (p Producer) Name() string {
	return p.name
}

func (p *Producer) SetExchangeBinds(eb []*ExchangeBinds) *Producer {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	if p.state != StateOpened {
		p.exchangeBinds = eb
	}
	return p
}

func (p *Producer) Open() error {
	if p.mq == nil {
		return errors.New("MQ: Bad producer")
	}

	// Open期间不允许对channel做任何操作
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if len(p.exchangeBinds) <= 0 {
		return errors.New("MQ: No exchangeBinds found. You should SetExchangeBinds before open.")
	}

	if p.state == StateOpened {
		return errors.New("MQ: Producer had been opened")
	}

	ch, err := p.mq.channel()
	if err != nil {
		return fmt.Errorf("MQ: Create channel failed, %v", err)
	}

	if err = p.applyExchangeBinds(ch); err != nil {
		ch.Close()
		return err
	}

	p.ch = ch
	p.state = StateOpened
	p.closeC = make(chan *amqp.Error, 1)
	p.ch.NotifyClose(p.closeC)

	go p.keepalive()
	return nil
}

func (p *Producer) Close() {
	if p.stopC != nil {
		close(p.stopC)
		p.stopC = nil
	}
}

func (p *Producer) Publish(exchange, routeKey string, msg *PublishMsg) error {
	if msg == nil {
		return errors.New("MQ: Nil publish msg")
	}

	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if p.state != StateOpened {
		return fmt.Errorf("MQ: Producer unopened, now state is %d", p.state)
	}
	pub := amqp.Publishing{
		ContentType:     msg.ContentType,
		ContentEncoding: msg.ContentEncoding,
		DeliveryMode:    msg.DeliveryMode,
		Priority:        msg.Priority,
		Timestamp:       msg.Timestamp,
		Body:            msg.Body,
	}
	return p.ch.Publish(exchange, routeKey, false, false, pub)
}

func (p *Producer) State() uint8 {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.state
}

func (p *Producer) keepalive() {
	select {
	case <-p.stopC:
		// 正常关闭
		log.Warn("MQ: Producer(%s) shutdown normally.", p.name)
		p.mutex.Lock()
		p.ch.Close()
		p.ch = nil
		p.state = StateClosed
		p.mutex.Unlock()

	case err := <-p.closeC:
		if err == nil {
			log.Error("MQ: Producer(%s)'s channel was closed, but Error detail is nil", p.name)
		} else {
			log.Error("MQ: Producer(%s)'s channel was closed, code:%d, reason:%s", p.name, err.Code, err.Reason)
		}

		// channel被异常关闭了
		p.mutex.Lock()
		p.state = StateReopening
		p.mutex.Unlock()

		maxRetry := 99999999
		for i := 0; i < maxRetry; i++ {
			time.Sleep(time.Second)
			if p.mq.State() != StateOpened {
				log.Warn("MQ: Producer(%s) try to recover channel for %d times, but mq's state != StateOpened", p.name, i+1)
				continue
			}
			if e := p.Open(); e != nil {
				log.Error("MQ: Producer(%s) recover channel failed for %d times, Err:%v", p.name, i+1, e)
				continue
			}
			log.Info("MQ: Producer(%s) recover channel OK. Total try %d times", p.name, i+1)
			return
		}
		log.Error("MQ: Producer(%s) try to recover channel over maxRetry(%d), so exit", p.name, maxRetry)
	}
}

func (p *Producer) applyExchangeBinds(ch *amqp.Channel) (err error) {
	if ch == nil {
		return errors.New("MQ: Nil producer channel")
	}

	for _, eb := range p.exchangeBinds {
		if eb.Exch == nil {
			return errors.New("MQ: Nil exchange found.")
		}
		if len(eb.Bindings) <= 0 {
			return fmt.Errorf("MQ: No bindings queue found for exchange(%s)", eb.Exch.Name)
		}
		// declare exchange
		ex := eb.Exch
		if err = ch.ExchangeDeclare(ex.Name, ex.Kind, ex.Durable, ex.AutoDelete, ex.Internal, ex.NoWait, ex.Args); err != nil {
			return fmt.Errorf("MQ: Declare exchange(%s) failed, %v", ex.Name, err)
		}
		// declare and bind queues
		for _, b := range eb.Bindings {
			if b == nil {
				return fmt.Errorf("MQ: Nil binding found, exchange:%s", ex.Name)
			}
			if len(b.Queues) <= 0 {
				return fmt.Errorf("MQ: No queues found for exchange(%s)", ex.Name)
			}
			for _, q := range b.Queues {
				if q == nil {
					return fmt.Errorf("MQ: Nil queue found, exchange:%s", ex.Name)
				}
				if _, err = ch.QueueDeclare(q.Name, q.Durable, q.AutoDelete, q.Exclusive, q.NoWait, q.Args); err != nil {
					return fmt.Errorf("MQ: Declare queue(%s) failed, %v", q.Name, err)
				}
				if err = ch.QueueBind(q.Name, b.RouteKey, ex.Name, b.NoWait, b.Args); err != nil {
					return fmt.Errorf("MQ: Bind exchange(%s) <--> queue(%s) failed, %v", ex.Name, q.Name, err)
				}
			}
		}
	}
	return nil
}
