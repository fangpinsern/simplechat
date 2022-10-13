package handlers

import (
	"context"
	"gochat/services/chat"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var (
	// writeWait = 10 * time.Second
	maxMessageSize int64 = 1024
	pongWait = 60 * time.Second
	// pingPeriod = (pongWait * 9) / 10
	// defaultBroadcastQueueSize = 10000
)

func ServeWS(ctx context.Context, c *chat.ChatInstance) func (w http.ResponseWriter, r *http.Request) {
	handler := func (w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		ws, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer ws.Close()

		// add authentication

		keys, ok := r.URL.Query()["user"]
		if !ok || len(keys[0]) < 1 {
			log.Println("Url Param 'user' is missing")
			return
		}
		log.Println(keys[0])
		userId := keys[0]

		sess := c.NewSession(ws)
		close := c.Bind(userId, sess.GetSessionId())
		defer close()

		ws.SetReadLimit(maxMessageSize)
		ws.SetReadDeadline(time.Now().Add(60*time.Second))
		ws.SetPongHandler(func(string) error {
			ws.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})
		go chat.Ping(ws)
		
		for {
			var msg chat.Message
			err := ws.ReadJSON(&msg)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err) {
					// split the different type of closes
					log.Println("error occured, unexpected close ", err)
				}
				break
			}
			msg.Sender = userId
			c.Broadcast <- msg
		}
	}

	return handler
}