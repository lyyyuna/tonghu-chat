package agent

import (
	"github.com/gorilla/websocket"
	"github.com/lyyyuna/tonghu-chat/pkg/broker"
)

// One agent per connection
type Agent struct {
	cb   *broker.ChatBroker
	conn *websocket.Conn
}

// NewAgent creates new connection agent instance
func NewAgent(conn *websocket.Conn, cb *broker.ChatBroker) *Agent {
	return &Agent{
		cb:   cb,
		conn: conn,
	}
}

func (a *Agent) HandleConn(conn *websocket.Conn, req *initConReq) {
	a.conn.SetCloseHandler(func(code int, text string) error {
		return nil
	})
}
