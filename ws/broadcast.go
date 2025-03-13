package ws

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"sync"
)

type Broadcast struct {
	subscribers map[string]chan []byte
	lock        sync.Mutex
}

func NewBroadcast() *Broadcast {
	return &Broadcast{
		subscribers: make(map[string]chan []byte),
	}
}

func (s *Broadcast) Subscribe() (string, chan []byte) {
	s.lock.Lock()
	defer s.lock.Unlock()

	subscriberId := genSubscriberId()
	log.Println("new subscriber", subscriberId)
	ch := make(chan []byte, 1)
	s.subscribers[subscriberId] = ch
	return subscriberId, ch
}

func (s *Broadcast) Unsubscribe(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.subscribers, key)
}

func (s *Broadcast) Broadcast(msg []byte) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, ch := range s.subscribers {
		go func() {
			ch <- msg
		}()
	}
}

func genSubscriberId() string {
	bytes := make([]byte, 32)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
