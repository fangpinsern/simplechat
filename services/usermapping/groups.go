package usermapping

import "context"

type InMemoryGroupCache struct {
	userToGroupMap map[string]map[string]bool // for closing connection
	groupToUserMap map[string]map[string]bool // for sending messages
}

func NewInMemoryGroupCache(ctx context.Context) *InMemoryGroupCache {
	return &InMemoryGroupCache{
		userToGroupMap: make(map[string]map[string]bool),
		groupToUserMap: make(map[string]map[string]bool),
	}
}

func (t *InMemoryGroupCache) Add(user, group string) error {
	_, ok := t.userToGroupMap[user]
	if !ok {
		t.userToGroupMap[user] = make(map[string]bool)
	}

	t.userToGroupMap[user][group] = true

	_, ok = t.groupToUserMap[group]
	if !ok {
		t.userToGroupMap[group] = make(map[string]bool)
	}

	t.userToGroupMap[group][user] = true
	return nil
}

func (t *InMemoryGroupCache) Delete(user string, fn func(string)) error {
	// delete user
	userGroups := t.userToGroupMap[user]
	
	for groupId := range userGroups {
		delete(t.groupToUserMap[groupId], user) 
		// notification to group that is deleted
		fn(groupId)
	}

	delete(t.userToGroupMap, user)

	return nil
}

func (t *InMemoryGroupCache) GetGroups(user string) []string {
	keys := make([]string, 0)
	groups := t.userToGroupMap[user]
	for groupId := range groups {
		keys = append(keys, groupId)
	}
	return keys
}

func (t *InMemoryGroupCache) GetUsers(group string) []string {
	keys := make([]string, 0)
	users := t.userToGroupMap[group]
	for userId := range users {
		keys = append(keys, userId)
	}
	return keys
}

