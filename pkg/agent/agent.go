package agent

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/lyyyuna/tonghu-chat/pkg/chat"
	"go.uber.org/zap"
	"io"
	"time"
)

// One agent per connection
type Agent struct {
	cb       chat.ChatBroker
	store    chat.ChatStore
	conn     *websocket.Conn
	channel  *chat.Channel
	user     *chat.User
	done     chan struct{}
	closed   bool
	closeSub func()
}

// NewAgent creates new connection agent instance
func NewAgent(conn *websocket.Conn, cb chat.ChatBroker, store chat.ChatStore) *Agent {
	return &Agent{
		cb:    cb,
		conn:  conn,
		store: store,
	}
}

func (a *Agent) HandleConn(conn *websocket.Conn, req *initConReq) {
	a.conn.SetCloseHandler(func(code int, text string) error {
		a.closed = true
		a.done <- struct{}{}
		return nil
	})

	ch, err := a.store.GetChannel(req.Channel)
	if err != nil {
		writeFatal(a.conn, fmt.Sprintf("agent: unable to find chat: %v", err))
	}

	user, err := ch.Join(req.UID)
	if err != nil {
		writeFatal(a.conn, fmt.Sprintf("agent: unable to join chat: %v", err))
		return
	}
	a.user = user

	mc := make(chan *chat.Message)
	var closer func()
	closer, err = a.cb.SubscribeNew(req.Channel, user.UID, mc)

	if err != nil {
		writeFatal(a.conn, fmt.Sprintf("agent: unable to subscribe to chat updates due to: %v. closing connection", err))
		return
	}
	a.closeSub = closer

}

func (a *Agent) loop(mc chan *chat.Message) {
	go func() {
		for {
			if a.closed {
				return
			}
			_, r, err := a.conn.NextReader()
			if err != nil {
				writeErr(a.conn, err.Error())
			}
			a.handleClientMsg(r)
		}
	}()

	go func() {
		defer a.closeSub()
		defer a.conn.Close()
		for {
			select {
			case m := <-mc:
				a.conn.WriteJSON(msg{
					Type: chatMsg,
					Data: m,
				})
			case <-a.done:
				return
			}
		}
	}()
}

func (a *Agent) handleClientMsg(r io.Reader) {
	var message struct {
		Type msgT            `json:"type"`
		Data json.RawMessage `json:"data,omitempty"`
	}

	err := json.NewDecoder(r).Decode(&message)
	if err != nil {
		writeErr(a.conn, fmt.Sprintf("invalid message format: %v", err))
		return
	}

	switch message.Type {
	case chatMsg:
		a.handleChatMsg(message.Data)
	default:
		writeErr(a.conn, fmt.Sprintf("unknown message format"))
	}
}

func (a *Agent) handleChatMsg(raw json.RawMessage) {
	var textMessage struct {
		Seq  uint64 `json:"seq"`
		Text string `json:"text"`
	}

	err := json.Unmarshal(raw, &textMessage)
	if err != nil {
		writeErr(a.conn, fmt.Sprintf("invalid message format: %v", err))
		return
	}

	if textMessage.Text == "" {
		writeErr(a.conn, fmt.Sprintf("empty message"))
		return
	}

	if len(textMessage.Text) > 1024 {
		writeErr(a.conn, fmt.Sprintf("message too long"))
		return
	}

	// send to chat broker
	err = a.cb.Send(a.channel.Name, &chat.Message{
		Meta:     nil,
		Time:     time.Now().UnixNano(),
		Seq:      textMessage.Seq,
		Text:     textMessage.Text,
		FromUID:  a.user.UID,
		FromName: a.user.DisplayName,
	})

	if err != nil {
		writeErr(a.conn, fmt.Sprintf("could not forward your message"))
		zap.S().Infof("Could not forward your message, err: %v", err)
		return
	}
}
