package wetalk

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
)

type MessageHandler func(msg *Message) *Message

type Client struct {
	Conn         *websocket.Conn
	ping         int64
	eventMapping map[string]MessageHandler
	Request      *http.Request
	UID          string
}

func (this *Client) On(event string, handler MessageHandler) {
	this.eventMapping[event] = handler
}

func (this *Client) OnClose(cb func(code int, text string) error) {
	this.Conn.SetCloseHandler(func(code int, text string) error {
		return cb(code, text)
	})
}

func (this *Client) Send(event string, data JSON) {
	bytes, _ := json.Marshal(data)
	this.Conn.WriteMessage(websocket.TextMessage, bytes)
}

func (this *Client) Reply(msg *Message) {
	bytes, _ := json.Marshal(msg)
	this.Conn.WriteMessage(websocket.TextMessage, bytes)
}
