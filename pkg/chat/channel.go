package chat

import (
	"errors"
)

type Channel struct {
	Name    string           `json:"name"`
	Members map[string]*User `json:"members"`
}

func NewChannel(name string) *Channel {
	return &Channel{
		Name:    name,
		Members: make(map[string]*User),
	}
}

// Chat errors
var (
	errAlreadyRegistered = errors.New("chat: uid already registered in this chat")
	errNotRegistered     = errors.New("chat: not a member of this channel")
	errInvalidSecret     = errors.New("chat: invalid secret")
)

// Register registers user with a chat and returns secret which should
// be stored on the client side, and used for subsequent join requests
func (c *Channel) Register(u *User) error {
	if _, ok := c.Members[u.UID]; ok {
		return errAlreadyRegistered
	}

	c.Members[u.UID] = u
	return nil
}

// Join attempts to join user to chat
func (c *Channel) Join(uid, secret string) (*User, error) {
	u, ok := c.Members[uid]
	if !ok {
		return nil, errNotRegistered
	}
	if u.Secret != secret {
		return nil, errInvalidSecret
	}

	return u, nil
}

func (c *Channel) Leave(uid string) {
	delete(c.Members, uid)
}

func (c *Channel) ListMembers(uid string) []*User {
	if len(c.Members) < 1 {
		return nil
	}

	var members []*User
	for _, v := range c.Members {
		members = append(members, v)
	}

	return members
}
