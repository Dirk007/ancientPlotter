package broadcast

import (
	"context"
	"sync"
)

type BroadcastChannel[T any] struct {
	mux      sync.Mutex
	channels map[ReceiverID]chan<- T
}

func NewBroadcastChannel[T any]() *BroadcastChannel[T] {
	return &BroadcastChannel[T]{
		channels: make(map[ReceiverID]chan<- T),
	}
}

func (b *BroadcastChannel[T]) newUniqueID() ReceiverID {
	var newID ReceiverID
	for {
		newID = NewReceiverID()
		_, exists := b.channels[newID]
		if !exists {
			break
		}
	}
	return newID
}

func (b *BroadcastChannel[T]) Register() (ReceiverID, <-chan T) {
	b.mux.Lock()
	defer b.mux.Unlock()

	newChannel := make(chan T)
	newID := b.newUniqueID()

	b.channels[newID] = newChannel
	return newID, newChannel
}

func (b *BroadcastChannel[T]) Remove(id ReceiverID) {
	b.mux.Lock()
	defer b.mux.Unlock()

	ch, ok := b.channels[id]
	if !ok {
		return
	}
	close(ch)
	delete(b.channels, id)
}

func (b *BroadcastChannel[T]) Broadcast(ctx context.Context, content T) error {
	b.mux.Lock()
	defer b.mux.Unlock()

	for _, ch := range b.channels {
		select {
		case ch <- content:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}
