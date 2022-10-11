package chat

import (
	"context"
	"gochat/services/session"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Text string `json:"data"`
	Id string `json:"id"`
	Type string `json:"type"`
	Timestamp int64 `json:"timestamp"`
	Sender string `json:"sender"`
	Receiver string `json:"receiver"`
}

type ChatInstance struct {
	broadcast chan Message
	quit chan struct{}
	sessions *session.Sessions
	sessionsLookup *session.UserMapping

	groups *Groups
	groupsLookup *UserGroupMapping

}

var upgrader = websocket.Upgrader{}
var (
	writeWait = 10 * time.Second
	maxMessageSize int64 = 1024
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	defaultBroadcastQueueSize = 10000
)

const MessageTypeMessage = "message"

func NewChatInstance(ctx context.Context) *ChatInstance {
	chat := ChatInstance{
		broadcast: make(chan Message, defaultBroadcastQueueSize),
		quit: make(chan struct{}),
		sessions: session.NewSessionsInstance(ctx),
		sessionsLookup: session.NewUserMappingInstance(ctx),
		groups: NewGroupsInstance(ctx),
		groupsLookup: NewUserGroupMappingInstance(ctx),
	}

	go chat.worker()
	return &chat
}

func (c *ChatInstance) worker() {
	// getStatus := func(userId string) string {
	// 	sessions := c.
	// }
	loop:
	for {
		select {
		case <- c.quit:
			log.Println("Quit")
			break loop
		case msg, ok := <- c.broadcast:
			if !ok {
				break loop
			}

			log.Println("processing message: ", msg)

			// switch msg.Type {
			// case MessageTypeMessage:
			// 	// do some db thing
			// default:
			// }
			c.Broadcast(&msg)
		}

	}
}

func (c *ChatInstance) Broadcast(message *Message) error {
	receiver := message.Receiver
	if receiver == "" {
		log.Println("WARN: receiver is empty")
	}

	sessions := c.sessionsLookup.GetSessionIdsOfUser(receiver)

	for _, sid := range sessions {
		sess := c.sessions.GetSession(sid)
		if sess == nil {
			continue
		}

		err := sess.GetConn().WriteJSON(message)
		if err != nil {
			c.Clear(sess)
			return err
		}
	}

	return nil

}

func (c *ChatInstance) newSession(ws *websocket.Conn) *session.Session {
	sess := session.NewSession(ws)
	c.sessions.InsertSession(sess)
	return sess
}

func (c *ChatInstance) Bind(userId, sessionId string) func() {
	c.sessionsLookup.Add(userId,sessionId)
	return func() {
		// function to stop the session
		session := c.sessions.GetSession(sessionId)
		c.Clear(session)
	}
}

func (c *ChatInstance) Clear(sess *session.Session) {
	if sess == nil {
		return
	}

	sess.GetConn().Close()
	sessId := sess.GetSessionId()

	c.sessionsLookup.DeleteSession(sessId)
	c.sessions.DeleteSession(sessId)
}

func (c *ChatInstance) ServeWS() func (w http.ResponseWriter, r *http.Request) {
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

		userId := "helloworld"

		sess := c.newSession(ws)
		close := c.Bind(userId, sess.GetSessionId())
		defer close()

		ws.SetReadLimit(maxMessageSize)
		ws.SetReadDeadline(time.Now().Add(60*time.Second))
		ws.SetPongHandler(func(string) error {
			ws.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})
		go ping(ws)
		
		for {
			var msg Message
			err := ws.ReadJSON(&msg)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err) {
					// split the different type of closes
					log.Println("error occured, unexpected close ", err)
				}
				break
			}
			msg.Sender = userId
			c.broadcast <- msg
		}
	}

	return handler
}

func ping(ws *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		ws.Close()
	}()
	loop:
	for {
		<-ticker.C
			log.Println("ponging from here")
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				// Don't use return, it will not trigger the defer function.
				break loop
			}
	}
}