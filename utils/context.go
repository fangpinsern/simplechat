package utils

import (
	"context"
	"log"

	"github.com/golang-jwt/jwt/v4"
)

const (
	KEY_USER_INFO = "user_info"
	KEY_USERNAME = "username"
	KEY_USER_ID = "id"
)

func GetUserId(ctx context.Context) string {
	token := ctx.Value("user").(*jwt.Token).Claims.(jwt.MapClaims)
	if token[KEY_USER_INFO] == nil {
		log.Println("token does not have user info?")
		return ""
	}
	userinfo := token[KEY_USER_INFO].(map[string]interface{})
	
	userId := userinfo[KEY_USER_ID]
	if userId == nil {
		return ""
	}
	return userId.(string)
}