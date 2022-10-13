package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AuthInstance struct {
	userdata map[string]User
}



var jwtKey = []byte("my_secret_key")

var users = map[string]User{
	"user1": {
		Username: "user1",
		Password: "password1",
	},
	"user2": {
		Username: "user2",
		Password: "password2",
	},
}

// Create a struct to read the username and password from the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewAuthorizeInstance(ctx context.Context, userDb map[string]User) *AuthInstance {
	return &AuthInstance{
		userdata: userDb,
	}
}

func (a *AuthInstance) Login(cred Credentials) (string,error)  {
	userInfo, ok := users[cred.Username]

	if !ok || userInfo.Password != cred.Password {
		return "", jwt.ErrInvalidKey
	}

	expirationTime := time.Now().Add(5 * time.Hour)
	claims := &Claims{
		Username: cred.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: expirationTime,
			},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func validateToken(tokenString string) (*Claims,error) {
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid || err == jwt.ErrTokenExpired {
			return claims, err
		}
		return claims, err
	}

	if !tkn.Valid {
		return claims, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

func (a *AuthInstance) ValidateToken(tokenString string) (*Claims,error) {
	return validateToken(tokenString)
}

// func (a *AuthInstance) Refresh(tokenString string) (string,error) {
// 	claims, err := validateToken(tokenString)
// 	if err != nil {
// 		return "", err
// 	}

// 	// should have a refresh token?
// }