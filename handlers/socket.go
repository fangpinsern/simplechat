package handlers

import (
	"context"
	"gochat/services/chat"
	"gochat/utils"
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

		userId:= utils.GetUserId(r.Context())

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
			msg.Timestamp = time.Now().Unix()
			log.Println("message is ", userId)
			c.Broadcast <- msg
		}
	}

	return handler
}