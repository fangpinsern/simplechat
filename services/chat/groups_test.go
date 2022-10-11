package chat_test

import (
	"context"
	"gochat/services/chat"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGroups(t *testing.T) {
	t.Parallel()

	groups := chat.NewGroupsInstance(context.Background())

	groupName := "testGroup"
	groupTypeTeam := "TEAM"
	admin := "ADMIN_USER"
	group := chat.NewGroup(groupTypeTeam, groupName, admin)

	groupId := group.GetGroupId()
	require.NotEmpty(t, groupId)

	require.Equal(t, groupName, group.GetGroupName())
	require.Equal(t, groupTypeTeam, group.GetGroupType())
	require.Equal(t, []string{admin}, group.GetAdmins())

	groups.InsertGroup(group)
	retGroup := groups.GetGroup(groupId)
	require.Equal(t, group, retGroup)

	groups.DeleteGroup(groupId)
	emptyGroups := groups.GetGroup(groupId)
	require.Empty(t, emptyGroups)
}