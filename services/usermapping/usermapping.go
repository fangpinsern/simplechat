package usermapping

import (
	"context"
	"sync"
)

// one user can have many sessions
// one session can only have 1 user
type InMemorySessionMap struct {
	mu sync.RWMutex
	sessionToUserMap map[string]string
	userToSessionMap map[string]map[string]bool
}

func NewInMemorySessionToUserMap(ctx context.Context) *InMemorySessionMap {
	return &InMemorySessionMap{
		sessionToUserMap: make(map[string]string),
		userToSessionMap: make(map[string]map[string]bool),
	}
}

func (t *InMemorySessionMap) GetUserSessionIds(userId string) []string {
	t.mu.RLock()
	sessionIds := t.userToSessionMap[userId]
	if sessionIds == nil {
		return []string{}
	}

	returnArr := make([]string, 0)

	for sessionId := range sessionIds {
		returnArr = append(returnArr, sessionId)
	}

	t.mu.RUnlock()
	return returnArr
}

func (t *InMemorySessionMap) GetSessionUserId(sessionId string) []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	userId, ok := t.sessionToUserMap[sessionId]
	if !ok {
		return []string{}
	}

	return []string{userId}
}

func (t *InMemorySessionMap) Add(userId, sessionId string) {
	t.mu.Lock()
	_, ok := t.userToSessionMap[userId]
	if !ok {
		t.userToSessionMap[userId] = make(map[string]bool)
	}
	t.userToSessionMap[userId][sessionId] = true

	t.sessionToUserMap[sessionId] = userId
	
	t.mu.Unlock()
}

func (t *InMemorySessionMap) DeleteSession(sessionId string) {
	t.mu.Lock()
	userId := t.sessionToUserMap[sessionId]

	delete(t.userToSessionMap[userId], sessionId)

	delete(t.sessionToUserMap, sessionId)
	
	t.mu.Unlock()
}

func (t *InMemorySessionMap) DeleteUser(userId string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	sessions := t.userToSessionMap[userId]
	for sessionId := range sessions {
		delete(t.sessionToUserMap, sessionId)
	}

	delete(t.userToSessionMap, userId)
}







// func (s *InMemorySocketMap) BindUser(user string, conn *websocket.Conn) error {
// 	userList := s.Users
// 	userList[user] = conn
// 	return nil
// }

// func (s *InMemorySocketMap) UnbindUser(user string) error {
// 	userList := s.Users

// 	delete(userList, user)
// 	return nil
// }

// func (s *InMemorySocketMap) GetUserSock(user string) (*websocket.Conn, error) {
// 	userList := s.Users
// 	conn, ok := userList[user]
// 	if !ok {
// 		err := fmt.Errorf("user connection does not exist. routed to notification center")
// 		return nil, err
// 	}
// 	return conn, nil
// }