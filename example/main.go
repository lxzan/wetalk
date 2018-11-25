package main

import (
	"im/ws"
	"net/http"
)

type Form map[string]string

func main() {
	ws.AuthHandler = func(r *http.Request) bool {
		r.Header.Set("X-User-ID", "1234")
		return true
	}

	socket := ws.NewServer()

	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		socket.OnConnect(writer, request, Onconnect)
	})
	http.ListenAndServe(":3000", nil)

}

func Onconnect(client *ws.Client) {
	client.On("test", func(data string) ws.JSON {
		return ws.JSON{
			"test": "true",
		}
	})
}
