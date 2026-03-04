package notify

import "sync"

func NewBroker[T any]() *Broker[T] {
	return &Broker[T]{
		channels: make(map[chan T]struct{}),
	}
}

type Broker[T any] struct {
	m        sync.Mutex
	channels map[chan T]struct{}
}

func (b *Broker[T]) Register(c chan T) {
	b.m.Lock()
	defer b.m.Unlock()

	b.channels[c] = struct{}{}
}

func (b *Broker[T]) Deregister(c chan T) {
	b.m.Lock()
	defer b.m.Unlock()

	delete(b.channels, c)
}

func (b *Broker[T]) Send(data T) {
	b.m.Lock()
	defer b.m.Unlock()

	for c := range b.channels {
		c <- data
	}
}
