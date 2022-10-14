package auth_test

import (
	"context"
	"gochat/services/auth"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
)

var users = map[string]auth.User{
	"user1": {
		Id: "1",
		Username: "user1",
		Password: "password1",
	},
	"user2": {
		Id: "2",
		Username: "user2",
		Password: "password2",
	},
}

func TestAuth(t *testing.T) {
	t.Parallel()

	authService := auth.NewAuthorizeInstance(context.Background(), users)

	username := "user1"
	passwordCorrect := "password1"
	passwordWrong := "123"
	usernameWrong := "wrong"

	correctCreds := auth.Credentials{
		Username: username,
		Password: passwordCorrect,
	}

	// valid login
	loginToken, err := authService.Login(correctCreds)
	require.Nil(t, err)
	require.NotEmpty(t, loginToken)

	claims, err := authService.ValidateToken(loginToken)
	require.Nil(t, err)
	require.Equal(t, claims.UserInfo.Username, username)

	// invalid login - wrong password

	wrongpasswordCreds := auth.Credentials{
		Username: username,
		Password: passwordWrong,
	}

	loginToken, err = authService.Login(wrongpasswordCreds)
	require.Equal(t, jwt.ErrInvalidKey,err)
	require.Empty(t, loginToken)

	wrongUsernameCreds := auth.Credentials{
		Username: usernameWrong,
		Password: passwordCorrect,
	}

	loginToken, err = authService.Login(wrongUsernameCreds)
	require.Equal(t, jwt.ErrInvalidKey,err)
	require.Empty(t, loginToken)
}