package pubsub

import "sync"

// RWLock is a wrapper around [sync.RWLock] that manages subscribers that get notified
// after the write lock was released
type RWLock struct {
	lock sync.RWMutex

	// subs is nil before the first subscription, so that a zero [RWLock] is valid
	subs map[chan struct{}]struct{}
}

func (s *RWLock) Lock() {
	s.lock.Lock()
}
func (s *RWLock) Unlock() {
	s.lock.Unlock()

	s.publish()
}

func (s *RWLock) RLock() {
	s.lock.RLock()
}
func (s *RWLock) RUnlock() {
	s.lock.RUnlock()
}

func (s *RWLock) Subscribe() chan struct{} {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.subs == nil {
		s.subs = make(map[chan struct{}]struct{})
	}

	ch := make(chan struct{})

	s.subs[ch] = struct{}{}

	return ch
}

func (s *RWLock) Unsubscribe(ch chan struct{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.subs, ch)
}

func (s *RWLock) publish() {
	s.lock.RLock()
	defer s.lock.RUnlock()

	for ch := range s.subs {
		ch <- struct{}{}
	}
}
