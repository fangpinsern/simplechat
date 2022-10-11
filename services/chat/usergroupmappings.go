package chat

import (
	"context"
	"sync"
)

/**
UserGroupMapping maps user to group and vice versa
A user can join multiple groups
A group can have multiple users
*/

/*
mu - mutex to make it concurrency safe
userToGroupMap - maps users to group
groupToUserMap - maps group to users
*/
type UserGroupMapping struct {
	mu sync.RWMutex
	// can move this to DB when you have it
	userToGroupMap map[string]map[string]bool // many to many
	groupToUserMap map[string]map[string]bool // many to many
}

func NewUserGroupMappingInstance(ctx context.Context) *UserGroupMapping {
	return &UserGroupMapping{
		userToGroupMap: make(map[string]map[string]bool),
		groupToUserMap: make(map[string]map[string]bool),
	}
}

func (g *UserGroupMapping) Add(userId, groupId string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	_, ok := g.userToGroupMap[userId]
	if !ok {
		g.userToGroupMap[userId] = make(map[string]bool)
	}

	g.userToGroupMap[userId][groupId] = true
	
	_, ok = g.groupToUserMap[groupId]
	if !ok {
		g.groupToUserMap[groupId] = make(map[string]bool)
	}

	g.groupToUserMap[groupId][userId] = true

	return nil
}

/*
DeleteUserFromGroup deletes user from the group
fn - is a closure function that send a notification to the group when the message is sent
*/
func (g *UserGroupMapping) DeleteUserFromGroup(userId string, groupId string, fn func(string)) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	users := g.groupToUserMap[groupId]
	delete(users, groupId)

	groups := g.userToGroupMap[userId]
	delete(groups, groupId)

	fn(groupId)

	return nil
}

func (g *UserGroupMapping) GetGroups(userId string) []string {
	key := make([]string, 0)
	groups := g.userToGroupMap[userId]
	for groupId := range groups {
		key = append(key, groupId)
	}

	return key
}

func (g *UserGroupMapping) GetUsers(groupId string) []string {
	key := make([]string, 0)
	users := g.groupToUserMap[groupId]
	for userId := range users {
		key = append(key, userId)
	}

	return key
}



