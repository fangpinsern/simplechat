package chat

import (
	"context"
	"fmt"
	"log"

	"gochat/services/usermapping"
)

var chatInstance *ChatInstance

type ChatInstance struct{}

type Message struct {
	To string `json:"to"`
	From string `json:"from"`
	Text string `json:"text"`
}

func InitializeChatInstance(ctx context.Context) {
	chatInstance = NewChatInstance(ctx)
}

func GetChatInstance() *ChatInstance {
	return chatInstance
}

func NewChatInstance(ctx context.Context) *ChatInstance {
	return &ChatInstance{}
}

func (c *ChatInstance) SendMessage(msg Message) error {

	socketMap := usermapping.GetInMemorySocketMap()
	toConn, err := socketMap.GetUserSock(msg.To)

	if err != nil {
		log.Println("to user not connected: ", err)
		return fmt.Errorf("user not connected. sent to notification server")
	}

	if err != nil {
		log.Printf("error occured when marshalling json message: %s", err)
	}
	
	err = toConn.WriteJSON(msg)
	if err != nil {
		log.Println("Error during message writing:", err)
		return fmt.Errorf("error occured when sending message")
	}

	return nil

}

// database interactions are added in service
