package wetalk

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/lxzan/ts"
	"log"
	"net/http"
	"time"
)

var AuthHandler = func(r *http.Request) bool {
	return true
}

type JSON map[string]interface{}

type Message struct {
	Event string
	Id    string
}

var upgrader *websocket.Upgrader

type Server struct {
	Clients *ts.Map
}

var server *Server

type MessageHeader struct {
	ID    string `json:"_id"`
	Event string `json:"_event"`
	Ack   bool   `json:"_ack"`
}

func NewServer() *Server {
	server = &Server{
		Clients: ts.NewMap(),
	}

	upgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     AuthHandler,
	}
	return server
}

func (this *Server) OnConnect(w http.ResponseWriter, r *http.Request, cb func(client *Client)) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		conn:         conn,
		ping:         time.Now().Unix(),
		eventMapping: make(map[string]MessageHandler),
		Request:      r,
	}

	uid := r.Header.Get("X-User-Id")
	server.Clients.Set(uid, client)

	cb(client)
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if string(p) == "1" {
			client.conn.WriteMessage(websocket.TextMessage, []byte("2"))
			client.ping = time.Now().Unix()
		} else {
			var header = &MessageHeader{}
			json.Unmarshal(p, header)
			fn, exist := client.eventMapping[header.Event]
			if exist {
				if header.Ack {
					client.Reply(header, fn(string(p)))
				} else {
					fn(string(p))
				}
			}
		}
	}
}
