package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lyyyuna/tonghu-chat/pkg/chat"
	"net/http"
)

// WssServer represents websocket api service
type WssServer struct {
	upgrader websocket.Upgrader
	broker   chat.ChatBroker
	store    chat.ChatStore
}

// NewWssServer creates new websocket api
func NewWssServer(r *gin.Engine, br chat.ChatBroker, store chat.ChatStore) *WssServer {
	api := &WssServer{
		broker: br,
		store:  store,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
	r.GET("/connect", api.connect)
	return api
}

func (api *WssServer) connect(c *gin.Context) {
	w := c.Writer
	r := c.Request
	conn, err := api.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("error while upgrading to ws connection: %v", err), 200)
		return
	}
	req, err := api.waitConnInit(conn)
	if err != nil {
		if err == errConnClosed {
			return
		}
		writeFatal(conn, err.Error())
		return
	}

	if req.UID == "" || req.Channel == "" {
		writeFatal(conn, "wrong init request")
		return
	}

	agent := NewAgent(conn, api.broker, api.store)
	agent.HandleConn(conn, req)
}

type initConReq struct {
	Channel string  `json:"channel"`
	UID     string  `json:"uid"`
	Secret  string  `json:"secret"` // User secret
	LastSeq *uint64 `json:"last_seq"`
}

var errConnClosed = errors.New("connection closed")

func (api *WssServer) waitConnInit(conn *websocket.Conn) (*initConReq, error) {
	t, wsr, err := conn.NextReader()
	if err != nil || t == websocket.CloseMessage {
		return nil, errConnClosed
	}

	var req initConReq
	err = json.NewDecoder(wsr).Decode(&req)
	if err != nil {
		return nil, err
	}

	return &req, nil
}
