package chat

import (
	"context"
	"gochat/services/session"
	"log"
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
	Broadcast chan Message
	quit chan struct{}
	sessions *session.Sessions
	sessionsLookup *session.UserMapping

	groups *Groups
	groupsLookup *UserGroupMapping

}

var (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second	
	defaultBroadcastQueueSize = 10000
	pingPeriod = (pongWait * 9) / 10
)

const MessageTypeMessage = "message"

func NewChatInstance(ctx context.Context) *ChatInstance {
	chat := ChatInstance{
		Broadcast: make(chan Message, defaultBroadcastQueueSize),
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
		case msg, ok := <- c.Broadcast:
			if !ok {
				break loop
			}

			log.Println("processing message: ", msg)

			// switch msg.Type {
			// case MessageTypeMessage:
			// 	// do some db thing
			// default:
			// }
			c.BroadcastMessage(&msg)
		}

	}
}

func (c *ChatInstance) BroadcastMessage(message *Message) error {
	receiver := message.Receiver
	if receiver == "" {
		log.Println("WARN: receiver is empty")
	}

	group := c.groupsLookup.GetUsers(receiver)

	// if no group, assume the reciever is a userId. We create a new group
	if len(group) == 0 {
		// abstract this out to search for groups
		// look for matching group between sender and reciever
		newGroup := NewGroup("private", "", message.Sender)
		c.groups.InsertGroup(newGroup)
		c.groupsLookup.Add(message.Receiver, newGroup.GetGroupId())
		c.groupsLookup.Add(message.Sender, newGroup.GetGroupId())
		message.Receiver = newGroup.GetGroupId()
		group = c.groupsLookup.GetUsers(message.Receiver)
	}

	for _, userId := range group {
		if userId == message.Sender {
			continue
		}

		sessions := c.sessionsLookup.GetSessionIdsOfUser(userId)
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
	}
	return nil
}

func (c *ChatInstance) NewSession(ws *websocket.Conn) *session.Session {
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


func Ping(ws *websocket.Conn) {
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