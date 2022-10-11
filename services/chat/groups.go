package chat

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"sync"
)

type Group struct {
	id string
	groupType string
	groupName string
	admin map[string]bool // key is userid
}

type Groups struct {
	mu sync.RWMutex
	groups map[string]*Group // map groupId to group
}

func NewGroupsInstance(ctx context.Context) *Groups {
	return &Groups{
		groups: make(map[string]*Group),
	}
}

// TODO: Create type for better group creation
func NewGroup(groupType, groupName string, admin string) *Group {
	return &Group{
		id: randomString(32),
		groupName: groupName,
		groupType: groupType,
		admin: map[string]bool{admin: true},
	}
}

func (g *Group) GetGroupId() string {
	return g.id
}

func (g *Group) GetGroupName() string {
	return g.groupName
}

func (g *Group) GetAdmins() []string {
	adminList := make([]string, 0)
	for userId := range g.admin {
		adminList = append(adminList, userId)
	}
	return adminList
}

func (g *Group) GetGroupType() string {
	return g.groupType
}

func (g *Groups) InsertGroup(group *Group) {
	g.mu.Lock()
	g.groups[group.id] = group
	g.mu.Unlock()
}

func (g *Groups) GetGroup(groupId string) *Group {
	g.mu.RLock()
	group := g.groups[groupId]
	g.mu.RUnlock()
	return group
}

// TODO: pass context to get userId
func (g *Groups) DeleteGroup(groupId string) {
	g.mu.Lock()
	// have to check if user is admin
	delete(g.groups, groupId)
	g.mu.Unlock()
}

func randomString(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}