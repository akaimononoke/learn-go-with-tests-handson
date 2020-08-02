package poker

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type playerServerWebSocket struct {
	*websocket.Conn
}

func newPlayerServerWebSocket(w http.ResponseWriter, r *http.Request) *playerServerWebSocket {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed to upgrade connection to WebSockets: %v\n", err)
	}
	return &playerServerWebSocket{conn}
}

func (w *playerServerWebSocket) Write(p []byte) (n int, err error) {
	err = w.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (w *playerServerWebSocket) WaitForMessage() string {
	_, msg, err := w.ReadMessage()
	if err != nil {
		log.Printf("failed to read from websocket: %v\n", err)
	}
	return string(msg)
}
