package wetalk

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/lxzan/ts"
	"log"
	"net/http"
	"time"
)

type JSON map[string]interface{}

type Server struct {
	Clients   *ts.Map
	Upgrader  *websocket.Upgrader
	Conf      *ServerConfig
	OnConnect func(client *Client)
}

var server *Server

type ServerConfig struct {
	HeartbeatInterval time.Duration              // heartbeat interval, default 25s
	Resend            time.Duration              // resend message, default 5s
	Free              time.Duration              // free unused connections, default 30min
	Passport          func(r *http.Request) bool // pass or reject connection
}

func NewServer(conf *ServerConfig) *Server {
	server = &Server{
		Clients: ts.NewMap(),
		Conf:    conf,
		Upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     conf.Passport,
		},
	}
	return server
}

func (this *Server) ServeWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := this.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		Conn:         conn,
		ping:         time.Now().Unix(),
		eventMapping: make(map[string]MessageHandler),
		Request:      r,
	}

	uid := r.Header.Get("X-User-Id")
	server.Clients.Set(uid, client)
	this.OnConnect(client)

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if string(p) == "1" {
			client.Conn.WriteMessage(websocket.TextMessage, []byte("2"))
			client.ping = time.Now().Unix()
		} else {
			var msg = &Message{}
			json.Unmarshal(p, msg)
			fn, exist := client.eventMapping[msg.Header.Event]
			if exist {
				if msg.Header.Ack {
					client.Reply(fn(msg))
				} else {
					fn(msg)
				}
			}
		}
	}
}

//func (this *Server) OnConnect(cb func(client *Client)) {
//
//}
