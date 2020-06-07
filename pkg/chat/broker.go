package chat

import "io"

// ChatBroker represents chat broker interface
type ChatBroker interface {
	Subscribe(string, string, uint64, chan *Message) (io.Closer, error)
	SubscribeNew(string, string, chan *Message) (io.Closer, error)
	Send(string, *Message) error
}
