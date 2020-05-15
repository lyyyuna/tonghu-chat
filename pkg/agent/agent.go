package agent

import (
	"github.com/lyyyuna/tonghu-chat/pkg/chat"
)

// One agent per connection
type Agent struct {
	cb *ChatBroker
}

// ChatBroker represents chat broker interface
type ChatBroker interface {
	Subscribe(string, string, uint64, chan *chat.Message) (func(), error)
	SubscribeNew(string, string, chan *chat.Message) (func(), error)
	Send(string, *chat.Message) error
}

// NewAgent creates new connection agent instance
func NewAgent(cb *ChatBroker) *Agent {
	return &Agent{
		cb: cb,
	}
}