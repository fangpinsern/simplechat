package usermapping

import (
	"context"
	"fmt"

	"github.com/gorilla/websocket"
)

var socketMap *InMemorySocketMap

func GetInMemorySocketMap() *InMemorySocketMap {
	return socketMap
}

func InitializeSocketMap(ctx context.Context) {
	socketMap = NewInMemorySocketMap(ctx)
}

type InMemorySocketMap struct {
	Users map[string]*websocket.Conn
}

func NewInMemorySocketMap(ctx context.Context) *InMemorySocketMap {
	return &InMemorySocketMap{
		Users: make(map[string]*websocket.Conn),
	}
}

func (s *InMemorySocketMap) BindUser(user string, conn *websocket.Conn) error {
	userList := s.Users
	userList[user] = conn
	return nil
}

func (s *InMemorySocketMap) UnbindUser(user string) error {
	userList := s.Users

	delete(userList, user)
	return nil
}

func (s *InMemorySocketMap) GetUserSock(user string) (*websocket.Conn, error) {
	userList := s.Users
	conn, ok := userList[user]
	if !ok {
		err := fmt.Errorf("user connection does not exist. routed to notification center")
		return nil, err
	}
	return conn, nil
}