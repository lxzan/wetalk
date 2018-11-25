package main

import (
	"github.com/lxzan/wetalk"
	"net/http"
	"time"
)

func main() {
	socket := wetalk.NewServer(&wetalk.ServerConfig{
		HeartbeatInterval: 25 * time.Second,
		Resend:            5 * time.Second,
		Free:              15 * time.Minute,
		Passport: func(r *http.Request) bool {
			r.Header.Set("X-User-ID", "1234")
			return true
		},
	})

	socket.OnConnect = func(client *wetalk.Client) {
		client.On("test", func(msg *wetalk.Message) *wetalk.Message {
			return msg.Reply(wetalk.JSON{
				"greet": "Hello!",
			})
		})
	}

	http.HandleFunc("/ws", socket.ServeWebSocket)

	http.ListenAndServe(":3000", nil)

}

