package wetalk

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
)

type MessageHandler func(data string) JSON

type Client struct {
	conn         *websocket.Conn
	ping         int64
	eventMapping map[string]MessageHandler
	Request      *http.Request
	Id           string
}

func (this *Client) On(event string, handler MessageHandler) {
	this.eventMapping[event] = handler
}

func (this *Client) Send(event string, data JSON) {
	bytes, _ := json.Marshal(data)
	this.conn.WriteMessage(websocket.TextMessage, bytes)
}

func (this *Client) Reply(header *MessageHeader, data JSON) {
	data["_id"] = header.ID
	data["_event"] = header.Event
	bytes, _ := json.Marshal(data)
	this.conn.WriteMessage(websocket.TextMessage, bytes)
}
