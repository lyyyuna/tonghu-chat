package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

// API represents websocket api service
type API struct {
	upgrader websocket.Upgrader
}

// NewAPI creates new websocket api
func NewAPI(r *gin.Engine) *API {
	api := &API{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return false },
		},
	}
	r.GET("/connect", api.connect)
	return api
}

func (api *API) connect(c *gin.Context) {
	w := c.Writer
	r := c.Request
	conn, err := api.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("error while upgrading to ws connection: %v", err), 200)
	}
	req, err := api.waitConnInit(conn)
	if err != nil {
		if err == errConnClosed {
			return
		}
		writeErr(conn, err.Error())
		return
	}
}

type initConReq struct {
	Channel string  `json:"channel"`
	UID     string  `json:"uid"`
	Secret  string  `json:"secret"` // User secret
	LastSeq *uint64 `json:"last_seq"`
}

var errConnClosed = errors.New("connection closed")

func (api *API) waitConnInit(conn *websocket.Conn) (*initConReq, error) {
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
