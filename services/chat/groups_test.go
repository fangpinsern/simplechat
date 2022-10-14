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
	admin := "ADMIN_USER"
	group := chat.NewGroup(groupName, admin)

	groupId := group.GetGroupId()
	require.NotEmpty(t, groupId)

	require.Equal(t, groupName, group.GetGroupName())
	require.Equal(t, chat.GROUP_TYPE_GROUP, group.GetGroupType())
	require.Equal(t, []string{admin}, group.GetAdmins())

	groups.InsertGroup(group)
	retGroup := groups.GetGroup(groupId)
	require.Equal(t, group, retGroup)

	user1 := "a"
	user2 := "b"
	expectedGroupName := "ab"
	personalChatGroup := chat.NewPersonalChat(user1, user2)
	require.Equal(t, expectedGroupName, personalChatGroup.GetGroupId())
	require.Equal(t, chat.GROUP_TYPE_PRIVATE, personalChatGroup.GetGroupType())
	
	require.Contains(t, personalChatGroup.GetAdmins(), user1)
	require.Contains(t, personalChatGroup.GetAdmins(), user2)

	personalChatId := chat.MakePersonalChatGroupId(user1, user2)
	require.Equal(t, expectedGroupName, personalChatId)

	groups.DeleteGroup(groupId)
	emptyGroups := groups.GetGroup(groupId)
	require.Empty(t, emptyGroups)
}