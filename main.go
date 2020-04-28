package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer ws.Close()

		for {
			mt, message, err := ws.ReadMessage()
			if err != nil {
				break
			}
			if string(message) == "ping" {
				message = []byte("pong")
			}

			err = ws.WriteMessage(mt, message)
			if err != nil {
				break
			}
		}
	})
	r.Run()
}
