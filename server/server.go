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

	chatInstance := chat.NewChatInstance(context.Background())
	authInstance := auth.NewAuthorizeInstance(context.Background(), testUserDB)

	// mw := negroni
	myRouter.HandleFunc("/", NotImplemented).Methods("GET")
	myRouter.HandleFunc("/socket", handlers.ServeWS(context.Background(), chatInstance)).Methods("GET")
	myRouter.HandleFunc("/testpost", testpostHandler).Methods("POST")

	myRouter.HandleFunc("/auth/login", handlers.Login(context.Background(), authInstance)).Methods("POST")
	srv := &http.Server{
        Handler:      myRouter,
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

// func socketHandler(w http.ResponseWriter, r *http.Request) {
// 	keys, ok := r.URL.Query()["user"]
// 	if !ok || len(keys[0]) < 1 {
//         log.Println("Url Param 'user' is missing")
//         return
//     }
// 	log.Println(keys[0])
// 	userid := keys[0]

// 	conn,err := upgrader.Upgrade(w,r,nil)
// 	if err != nil {
// 		log.Print("Error during conenction upgradation: ", err)
// 		return
// 	}

// 	usermapping.GetInMemorySocketMap().BindUser(userid, conn)


// 	defer func(){
// 		usermapping.GetInMemorySocketMap().UnbindUser(userid)
// 		conn.Close()
// 	}()

// 	for {
// 		var msgJson chat.Message
// 		err := conn.ReadJSON(&msgJson)
// 		if err != nil {
// 			log.Println("error during message reading: ", err)
// 			break
// 		}

// 		log.Printf("Received: %s, To: %s", msgJson.Text, msgJson.To)

// 		toUser := msgJson.To
// 		fromUser := msgJson.From

// 		chatMessage := chat.Message{
// 			To: toUser,
// 			From: fromUser,
// 			Text: msgJson.Text,
// 		}
// 		response := chat.GetChatInstance().SendMessage(chatMessage)
// 		if !response.Success {
// 			errorMessage := chat.Message{
// 				To: fromUser,
// 				From: fromUser,
// 				Text: strconv.Itoa(response.Code),
// 			}
// 			chat.GetChatInstance().SendMessage(errorMessage)
// 		}
// 		if err != nil {
// 			log.Println("user does not exist: ", err)
// 			conn.WriteMessage(1, []byte("unsucessful"))
// 			continue
// 		}
// 	}
// }

// func home(w http.ResponseWriter, r *http.Request) {
//     fmt.Fprintf(w, "Index Page")
// }

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

