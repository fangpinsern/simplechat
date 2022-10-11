package main

import (
	"fmt"
	"gochat/server"

	_ "github.com/mattn/go-sqlite3"
)

// var database *sql.DB

func main() {

	// usermapping.InitializeSocketMap(context.Background())
	// chat.InitializeChatInstance(context.Background())

	// database1, err := sql.Open("sqlite3", "./helloworld.db")
	// database = database1
	// if err != nil {
	// 	log.Fatal("error ", err)
	// }

	// statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS chat (id INTEGER PRIMARY KEY, reciever TEXT, sender TEXT, message TEXT)")
	// if err != nil {
	// 	log.Fatal("error2 ", err)
	// }
	// statement.Exec()

	// rows, _ := database.Query("SELECT id, firstname, lastname FROM people")

	// var id int
	// var firstname string
	// var lastname string

	// for rows.Next() {
	// 	rows.Scan(&id, &firstname, &lastname)
	// 	fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
	// }

	fmt.Println("Helloworld")

	server.StartServer()
}