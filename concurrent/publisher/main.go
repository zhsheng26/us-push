package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

//
// @author zhangsheng
// @date 2019/7/31
//

type (
	subscriber chan interface{}
	topicFunc  func(v interface{}) bool
)

type Publisher struct {
	m           sync.RWMutex
	buffer      int
	timeout     time.Duration
	subscribers map[subscriber]topicFunc
}

func NewPublisher(publishTimeout time.Duration, buffer int) *Publisher {
	return &Publisher{
		buffer:      buffer,
		timeout:     publishTimeout,
		subscribers: make(map[subscriber]topicFunc),
	}
}

func (p *Publisher) Subscribe() chan interface{} {
	return p.SubscribeTopic(nil)
}

func (p *Publisher) SubscribeTopic(topic topicFunc) chan interface{} {
	ch := make(chan interface{}, p.buffer)
	p.m.Lock()
	p.subscribers[ch] = topic
	p.m.Unlock()
	return ch
}
func (p *Publisher) Evict(sub chan interface{}) {
	p.m.Lock()
	defer p.m.Unlock()
	delete(p.subscribers, sub)
	close(sub)
}
func (p *Publisher) Publish(v interface{}) {
	p.m.Lock()
	defer p.m.Unlock()
	var wg sync.WaitGroup
	for sub, topic := range p.subscribers {
		wg.Add(1)
		go p.sendTopic(sub, topic, v, &wg)
	}
	wg.Wait()
}

func (p Publisher) sendTopic(sub subscriber, topic topicFunc, v interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	if topic != nil && !topic(v) {
		return
	}
	select {
	case sub <- v:

	case <-time.After(p.timeout):

	}

}

func (p *Publisher) Close() {
	p.m.Lock()
	defer p.m.Unlock()
	for sub := range p.subscribers {
		delete(p.subscribers, sub)
		close(sub)
	}
}

func main() {
	p := NewPublisher(100*time.Millisecond, 10)
	defer p.Close()
	all := p.Subscribe()
	golang := p.SubscribeTopic(func(v interface{}) bool {
		if s, ok := v.(string); ok {
			return strings.Contains(s, "golang")
		}
		return false
	})
	p.Publish("hello,world")
	p.Publish("hello,golang")
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
	ch := make(chan int)
	go func() {
		for {
			select {
			case it := <-all:
				fmt.Println("all:", it)
			case <-ctx.Done():
				//ch <- 1
				close(ch)
				return
			}
		}
	}()
	go func() {
		for msg := range golang {
			fmt.Println("golang", msg)
		}
	}()
	<-ch
}
