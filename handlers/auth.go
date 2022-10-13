package handlers

import (
	"context"
	"encoding/json"
	"gochat/services/auth"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Username string `json:"username"`
	Token string `json:"token"`
}

func Login(ctx context.Context, authInstance *auth.AuthInstance) func (w http.ResponseWriter, r *http.Request) {
	handler := func (w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// get the body
		loginRequest := &LoginRequest{}
		err := json.NewDecoder(r.Body).Decode(loginRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		creds := auth.Credentials{
			Username: loginRequest.Username,
			Password: loginRequest.Password,
		}

		token, err := authInstance.Login(creds)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		res := LoginResponse{
			Username: loginRequest.Username,
			Token: token,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}

	return handler
}