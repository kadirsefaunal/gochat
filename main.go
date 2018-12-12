package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/googollee/go-socket.io"
)

// Message keeps message informations
type Message struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Message string `json:"message"`
}

// User keeps disconnected user information
type User struct {
	UserName string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	users := make(map[string]socketio.Socket)

	server.On("connection", func(so socketio.Socket) {
		so.On("set-user", func(user string) {
			users[user] = so
			fmt.Println(users)

			keys := make([]string, 0, len(users))
			for key := range users {
				keys = append(keys, key)
			}

			for _, socket := range users {
				socket.Emit("user-list", keys)
			}
		})

		so.On("message", func(msg string) {
			jsonMsg := Message{}
			json.Unmarshal([]byte(msg), &jsonMsg)
			fmt.Println(jsonMsg)

			toSo := users[jsonMsg.To]
			if toSo != nil {
				toSo.Emit("message", jsonMsg)
			}
		})

		so.On("disconnection", func() {
			var k string
			for key, socket := range users {
				if socket == so {
					k = key
					break
				}
			}

			delete(users, k)

			user := User{
				UserName: k,
				Message:  "User left!",
			}

			for _, socket := range users {
				socket.Emit("message", user)
			}

			fmt.Println(users)
			log.Println("Disconnect!")
		})
	})

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:5000 ...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
