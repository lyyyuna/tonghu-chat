package chat

// ChatBroker represents chat broker interface
type ChatBroker interface {
	Subscribe(string, string, uint64, chan *Message) (func(), error)
	SubscribeNew(string, string, chan *Message) (func(), error)
	Send(string, *Message) error
}
