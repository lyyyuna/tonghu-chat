package chat

type ChatStore interface {
	GetChannel(string) (*Channel, error)
}
