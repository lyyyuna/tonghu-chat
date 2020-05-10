package agent

import (
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

}

func (api *API) waitConnInit(conn *websocket.Conn) {
	t, wsr, err := conn.NextReader()
}
