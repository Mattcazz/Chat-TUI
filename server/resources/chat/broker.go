package chat

import (
	"sync"
	"time"

	"github.com/Mattcazz/Chat-TUI/pkg"
)

type Broker struct {
	mu          sync.RWMutex
	subscribers map[int64][]chan pkg.MsgResponse // convID -> listeners
}

func NewBroker() *Broker {
	return &Broker{
		subscribers: make(map[int64][]chan pkg.MsgResponse),
	}
}

func (b *Broker) Subscribe(convID int64) chan pkg.MsgResponse {
	ch := make(chan pkg.MsgResponse, 8)
	b.mu.Lock()
	b.subscribers[convID] = append(b.subscribers[convID], ch)
	b.mu.Unlock()
	return ch
}

func (b *Broker) Unsubscribe(convID int64, ch chan pkg.MsgResponse) {
	b.mu.Lock()
	defer b.mu.Unlock()
	subs := b.subscribers[convID]
	for i, sub := range subs {
		if sub == ch {
			b.subscribers[convID] = append(subs[:i], subs[i+1:]...) // remove the subscriber at index i
			break
		}
	}
}

func (b *Broker) Publish(convID int64, msg pkg.MsgResponse) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, ch := range b.subscribers[convID] {
		go func(c chan pkg.MsgResponse) {
			select {
			case c <- msg:
			case <-time.After(5 * time.Second): // client not consuming, drop message
			}
		}(ch)
	}
}
