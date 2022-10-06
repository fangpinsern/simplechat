package main

import (
	"context"
	"database/sql"
	"fmt"
	"gochat/services/usermapping"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

type Message struct {
	To string `json:"to"`
	From string `json:"from"`
	Text string `json:"text"`
}

var upgrader = websocket.Upgrader{}
var socketMap *usermapping.InMemorySocketMap


func socketHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["user"]
	if !ok || len(keys[0]) < 1 {
        log.Println("Url Param 'key' is missing")
        return
    }
	log.Println(keys[0])
	userid := keys[0]

	conn,err := upgrader.Upgrade(w,r,nil)
	if err != nil {
		log.Print("Error during conenction upgradation: ", err)
		return
	}

	socketMap.BindUser(userid, conn)

	defer func(){
		socketMap.UnbindUser(userid)
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
		toConn, err := socketMap.GetUserSock(toUser)
		if err != nil {
			log.Println("user does not exist: ", err)
			conn.WriteMessage(1, []byte("unsucessful"))
			continue
		}

		err = toConn.WriteJSON(msgJson)
		if err != nil {
			log.Println("Error during message writing:", err)
            break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Index Page")
}

func main() {

	socketMap = usermapping.NewInMemorySocketMap(context.Background())

	database, err := sql.Open("sqlite3", "./helloworld.db")
	if err != nil {
		log.Fatal("error ", err)
	}

	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	if err != nil {
		log.Fatal("error2 ", err)
	}
	statement.Exec()

	// statement, _ = database.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
	// statement.Exec("john", "doe")

	rows, _ := database.Query("SELECT id, firstname, lastname FROM people")

	var id int
	var firstname string
	var lastname string

	for rows.Next() {
		rows.Scan(&id, &firstname, &lastname)
		fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
	}

	fmt.Println("Helloworld")

	http.HandleFunc("/socket", socketHandler)
	http.HandleFunc("/", home)

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}