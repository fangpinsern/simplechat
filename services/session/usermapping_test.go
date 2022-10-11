package session_test

import (
	"context"
	"gochat/services/session"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserMapping(t *testing.T) {
	t.Parallel()

	userMap := session.NewUserMappingInstance(context.Background())

	// add user to user and session
	stubSessionId := "stubsession"
	stubUserId := "stubUse"

	userMap.Add(stubUserId, stubSessionId)
	userId := userMap.GetUserIdOfSession(stubSessionId)
	require.Equal(t, stubUserId, userId)

	sessionIds := userMap.GetSessionIdsOfUser(userId)
	require.Equal(t, []string{stubSessionId}, sessionIds)

	// session dont exist - userId is empty
	stubSessionDontExist := "stubSessionDontExist"
	userIdDontExist := userMap.GetUserIdOfSession(stubSessionDontExist)
	require.Equal(t, userIdDontExist, "")

	stubUserDontExist := "stubUserDontExist"
	emptySessionIds := userMap.GetSessionIdsOfUser(stubUserDontExist)
	require.Empty(t, emptySessionIds)

	// delete session
	userMap.DeleteSession(stubSessionId)
	sessionIds = userMap.GetSessionIdsOfUser(stubUserId)
	require.Empty(t, sessionIds)
	userId = userMap.GetUserIdOfSession(stubSessionId)
	require.Empty(t, userId)
}