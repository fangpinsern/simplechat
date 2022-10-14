package server

import (
	"context"
	"encoding/json"
	"gochat/handlers"
	"gochat/services/auth"
	"gochat/services/chat"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mfuentesg/go-jwtmiddleware"
	"github.com/urfave/negroni"
)

type Server struct {
}

// var upgrader = websocket.Upgrader{}
var(
	testUserDB = map[string]auth.User{
		"user1": {
			Username: "user1",
			Password: "password1",
		},
		"user2": {
			Username: "user2",
			Password: "password2",
		},
	}
)


func StartServer() {
	myRouter := mux.NewRouter().StrictSlash(true)
	authRouter := mux.NewRouter().PathPrefix("/api").Subrouter().StrictSlash(true)

	chatInstance := chat.NewChatInstance(context.Background())
	authInstance := auth.NewAuthorizeInstance(context.Background(), testUserDB)

	jwtmw := jwtmiddleware.New(jwtmiddleware.WithSignKey([]byte("my_secret_key")))

	myRouter.HandleFunc("/", NotImplemented).Methods("GET")
	myRouter.HandleFunc("/testpost", testpostHandler).Methods("POST")
	myRouter.HandleFunc("/auth/login", handlers.Login(context.Background(), authInstance)).Methods("POST")

	authRouter.HandleFunc("/ping", NotImplemented).Methods("GET")
	authRouter.HandleFunc("/socket", handlers.ServeWS(context.Background(), chatInstance)).Methods("GET")

	n := negroni.Classic()


	myRouter.PathPrefix("/api").Handler(n.With(
		negroni.HandlerFunc(jwtmw.HandlerNext),
		negroni.Wrap(authRouter),
	))


	n.UseHandler(myRouter)

	srv := &http.Server{
        Handler:      n,
        Addr:         "localhost:8080",
        // Good practice: enforce timeouts for servers you create!
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }

	log.Fatal(srv.ListenAndServe())
}

var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Not Implemented"))
  })

type testBody struct {
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
}

func testpostHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody testBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Println("error occured on decoding body: ", err)
		return
	}

	response, err := json.Marshal(&requestBody)
	if err != nil {
		log.Println("error occured on encoding body: ", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write(response)
}

