package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"gochat/services/chat"
	"gochat/services/usermapping"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

type Message struct {
	To string `json:"to"`
	From string `json:"from"`
	Text string `json:"text"`
}

var upgrader = websocket.Upgrader{}
var database *sql.DB


func socketHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["user"]
	if !ok || len(keys[0]) < 1 {
        log.Println("Url Param 'user' is missing")
        return
    }
	log.Println(keys[0])
	userid := keys[0]

	conn,err := upgrader.Upgrade(w,r,nil)
	if err != nil {
		log.Print("Error during conenction upgradation: ", err)
		return
	}

	usermapping.GetInMemorySocketMap().BindUser(userid, conn)


	defer func(){
		usermapping.GetInMemorySocketMap().UnbindUser(userid)
		conn.Close()
	}()

	for {
		var msgJson Message
		err := conn.ReadJSON(&msgJson)
		if err != nil {
			log.Println("error during message reading: ", err)
			break
		}

		log.Printf("Received: %s, To: %s", msgJson.Text, msgJson.To)

		toUser := msgJson.To
		fromUser := msgJson.From

		chatMessage := chat.Message{
			To: toUser,
			From: fromUser,
			Text: msgJson.Text,
		}
		err = chat.GetChatInstance().SendMessage(chatMessage)
		if err != nil {
			log.Println("user does not exist: ", err)
			conn.WriteMessage(1, []byte("unsucessful"))
			continue
		}

		// statement, err := database.Prepare("INSERT INTO chat (reciever, sender, message) VALUES (?, ?, ?)")
		// if err != nil {
		// 	log.Fatal("SOMETHING WENT WRONG", err)
		// 	return
		// }
		// _, err = statement.Exec(toUser, fromUser, msgJson.Text)
		// if err != nil {
		// 	log.Fatal(err)
		// }
	}
}

func home(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Index Page")
}

type testBody struct {
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
}

func testpostHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody testBody
	// json.Unmarshal(body, requestBody)
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

func startServer() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", home).Methods("GET")
	myRouter.HandleFunc("/socket", socketHandler).Methods("GET")
	myRouter.HandleFunc("/testpost", testpostHandler).Methods("POST")

	srv := &http.Server{
        Handler:      myRouter,
        Addr:         "localhost:8080",
        // Good practice: enforce timeouts for servers you create!
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }

	log.Fatal(srv.ListenAndServe())
}

func main() {

	usermapping.InitializeSocketMap(context.Background())

	database1, err := sql.Open("sqlite3", "./helloworld.db")
	database = database1
	if err != nil {
		log.Fatal("error ", err)
	}

	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS chat (id INTEGER PRIMARY KEY, reciever TEXT, sender TEXT, message TEXT)")
	if err != nil {
		log.Fatal("error2 ", err)
	}
	statement.Exec()

	rows, _ := database.Query("SELECT id, firstname, lastname FROM people")

	var id int
	var firstname string
	var lastname string

	for rows.Next() {
		rows.Scan(&id, &firstname, &lastname)
		fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
	}

	fmt.Println("Helloworld")

	startServer()
}