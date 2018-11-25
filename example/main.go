package main

import (
	"fmt"
	"github.com/lxzan/wetalk"
	"net/http"
	"time"
)

func main() {
	socket := wetalk.NewServer(&wetalk.ServerConfig{
		PingInterval: 25,
		Resend:       5,
		Free:         15 * time.Minute,
		Passport: func(r *http.Request) bool {
			r.Header.Set("UID", "caster")
			return true
		},
	})

	socket.OnConnect = func(client *wetalk.Client) {
		client.On("test", func(msg *wetalk.Message) *wetalk.Message {
			return msg.Reply(wetalk.JSON{
				"hello": "world",
			})
		})

		client.OnClose(func(code int, text string) error {
			println(fmt.Sprintf("%s disconnected.", client.UID))
			return nil
		})
	}

	http.HandleFunc("/socket", socket.ServeWebSocket)

	http.ListenAndServe(":3000", nil)

}
