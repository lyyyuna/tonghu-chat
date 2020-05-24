package store

import "github.com/lyyyuna/tonghu-chat/pkg/chat"

type ChatStore interface {
	GetChannel(string) (*chat.Channel, error)
}
