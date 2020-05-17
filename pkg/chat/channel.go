package chat

import ()

type Channel struct {
	Name    string           `json:"name"`
	Members map[string]*User `json:"members"`
}

func (c *Channel) Join(uid string) {

}
