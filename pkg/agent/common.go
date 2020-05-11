package agent

import "github.com/gorilla/websocket"

type msgT int

const (
	chatMsg msgT = iota
	historyMsg
	errorMsg
	infoMsg
	historyReqMsg
)

type msg struct {
	Type  msgT        `json:"type"`
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func writeErr(conn *websocket.Conn, err string) {
	conn.WriteJSON(msg{Error: err, Type: errorMsg})
}

func writeFatal(conn *websocket.Conn, err string) {
	conn.WriteJSON(msg{Error: err, Type: errorMsg})
	conn.Close()
}
