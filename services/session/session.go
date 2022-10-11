package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"sync"

	"github.com/gorilla/websocket"
)

type Session struct {
	id string
	conn *websocket.Conn
}

type Sessions struct {
	mu sync.RWMutex
	sessions map[string]*Session // map sessionId to session
}

func NewSessionsInstance(ctx context.Context) *Sessions {
	return &Sessions{
		sessions: make(map[string]*Session),
	}
}

func NewSession(conn *websocket.Conn) *Session {
	return &Session{
		id: randomString(32),
		conn: conn,
	}
}

func (s *Session) GetConn() *websocket.Conn {
	return s.conn
}

func (s *Session) GetSessionId() string {
	return s.id
}

func (s *Sessions) InsertSession(session *Session) {
	s.mu.Lock()
	s.sessions[session.id] = session
	s.mu.Unlock()
}

func (s *Sessions) GetSession(sessionId string) *Session {
	s.mu.RLock()
	session := s.sessions[sessionId]
	s.mu.RUnlock()
	return session
}

func (s *Sessions) DeleteSession(sessionId string) {
	s.mu.Lock()
	delete(s.sessions, sessionId)
	s.mu.Unlock()
}

func randomString(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}