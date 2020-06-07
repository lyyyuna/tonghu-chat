package chat

type ChatStore interface {
	GetChannel(string) (*Channel, error)
	SaveChannel(*Channel) error
}
