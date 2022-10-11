package session_test

import (
	"context"
	"gochat/services/session"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func TestSessions(t *testing.T) {
	t.Parallel()
	
	sessions := session.NewSessionsInstance(context.Background())

	stubWsConn := websocket.Conn{}
	session := session.NewSession(&stubWsConn)

	sessionId := session.GetSessionId()
	sessionConn := session.GetConn()
	require.Equal(t, &stubWsConn, sessionConn)

	sessions.InsertSession(session)

	retrievedSession := sessions.GetSession(sessionId)
	require.NotNil(t, retrievedSession)
	require.Equal(t, sessionId, retrievedSession.GetSessionId())
	require.Equal(t, sessionConn, retrievedSession.GetConn())

	sessions.DeleteSession(sessionId)

	retrievedSession = sessions.GetSession(sessionId)
	require.Nil(t, retrievedSession)
}