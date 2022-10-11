package usermapping

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
	server string // for extensiblity
}

type InMemorySessionTable struct {
	mu sync.RWMutex
	sessions map[string]*Session
}

func NewInMemorySessionTable(ctx context.Context) *InMemorySessionTable {
	return &InMemorySessionTable{
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

func randomString(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

func (s *InMemorySessionTable) InsertSession(sess *Session) {
	s.mu.Lock()
	s.sessions[sess.id] = sess
	s.mu.Unlock()
}

func (s *InMemorySessionTable) GetSession(id string) *Session {
	s.mu.RLock()
	sess := s.sessions[id]
	s.mu.RUnlock()
	return sess
}

func (s *InMemorySessionTable) DeleteSession(id string) {
	s.mu.Lock()
	delete(s.sessions, id)
	s.mu.Unlock()
}