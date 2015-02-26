package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ConnectHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Upgrading the socket connection")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Print("Connection upgraded")
	go func(c *websocket.Conn) {
		for {
			log.Print("Reading socket messages")
			if _, _, err := c.NextReader(); err != nil {
				log.Print(err)
				c.Close()
				break
			}
			log.Print("A message was processed")
		}
	}(conn)	
}

func BlockHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Blocking the socket connection")
	w.WriteHeader(http.StatusForbidden)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/ws/pass", ConnectHandler)
	r.HandleFunc("/ws/block", BlockHandler)

	err := http.ListenAndServe(":7082", r)
	if err != nil {
		log.Fatal(err)
	}
}
