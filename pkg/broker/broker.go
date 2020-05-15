package broker

import (
	"github.com/lyyyuna/tonghu-chat/pkg/chat"
)

// ChatBroker represents chat broker interface
type ChatBroker interface {
	Subscribe(string, string, uint64, chan *chat.Message) (func(), error)
	SubscribeNew(string, string, chan *chat.Message) (func(), error)
	Send(string, *chat.Message) error
}
