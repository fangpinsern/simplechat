package session

import (
	"context"
	"sync"
)

/**
Usermapping maps a user to a session
A user can have multiple sessions (mobile, computer etc)
But a session every session can only have 1 user
*/

/*
mu - mutex to make it concurrency safe
userToSessionMap - maps between userId and sessionId
*/
type UserMapping struct {
	mu sync.RWMutex
	userToSessionMap map[string]map[string]bool // one to many
	sessionToUserMap map[string]string // one to one
}

func NewUserMappingInstance(ctx context.Context) *UserMapping {
	return &UserMapping{
		userToSessionMap: make(map[string]map[string]bool),
		sessionToUserMap: make(map[string]string),
	}
}

func (u *UserMapping) GetSessionIdsOfUser(userId string) []string {
	u.mu.RLock()
	defer u.mu.RUnlock()
	sessionIds := u.userToSessionMap[userId]
	if sessionIds == nil {
		return []string{}
	}

	returnArr := make([]string, 0)

	for sessionId := range sessionIds {
		returnArr = append(returnArr, sessionId)
	}

	return returnArr
}

func (u *UserMapping) GetUserIdOfSession(sessionId string) string {
	u.mu.RLock()
	defer u.mu.RUnlock()

	userId, ok := u.sessionToUserMap[sessionId]
	if !ok {
		return ""
	}

	return userId
}

func (u *UserMapping) Add(userId, sessionId string) {
	u.mu.Lock()
	defer u.mu.Unlock()

	_, ok := u.userToSessionMap[userId]
	if !ok {
		u.userToSessionMap[userId] = make(map[string]bool)
	}

	u.userToSessionMap[userId][sessionId] = true
	u.sessionToUserMap[sessionId] = userId
}

func (u *UserMapping) DeleteSession(sessionId string) {
	u.mu.Lock()
	defer u.mu.Unlock()

	userId := u.sessionToUserMap[sessionId]

	delete(u.userToSessionMap[userId], sessionId)
	delete(u.sessionToUserMap, sessionId)
}

func (u *UserMapping) DeleteUser(userId string) {
	u.mu.Lock()
	defer u.mu.Unlock()

	sessions := u.userToSessionMap[userId]
	
	for sessionId := range sessions {
		delete(u.sessionToUserMap, sessionId)
	}

	delete(u.userToSessionMap, userId)
}