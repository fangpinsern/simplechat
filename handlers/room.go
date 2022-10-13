package handlers

import (
	"context"
	"gochat/services/chat"
	"net/http"
)

/*
GetGroups - gets the groups of the user
user has to be authorized to get the group
*/
func GetGroups(ctx context.Context, c *chat.ChatInstance) func (w http.ResponseWriter, r *http.Request) {
	handler := func (w http.ResponseWriter, r *http.Request) {
		
	}

	return handler
}

func CreateGroups(ctx context.Context, c *chat.ChatInstance) func (w http.ResponseWriter, r *http.Request) {
	handler := func (w http.ResponseWriter, r *http.Request) {
		
	}

	return handler
}