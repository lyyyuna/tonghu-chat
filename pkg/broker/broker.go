package broker

import (
	"io"
)

type Broker struct {
	mq MessageBroker
}

// MessageBroker represents a message broker interface
type MessageBroker interface {
	Send(string, []byte) error
	SubscribeSeq(string, string, uint64, func(uint64, []byte)) (io.Closer, error)
}

func NewBroker(mq MessageBroker) *Broker {
	return &Broker{
		mq: mq,
	}
}
