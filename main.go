package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	http.HandleFunc("/socket", socketHandler)
	http.HandleFunc("/", home)
	bind := "localhost:8099"
	log.Print(bind)
	log.Fatal(http.ListenAndServe(bind, nil))
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	var i int
	go func() {
		for {
			<-time.After(2 * time.Second)
			err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("ping %v", i)))
			if err != nil {
				log.Println("Error during message writing:", err)
			}
			i++
		}
	}()
	// The event loop
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			break
		}
		log.Printf("Received: %s", message)
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Error during message writing:", err)
			break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	index.ExecuteTemplate(w, "", nil)
}

var upgrader websocket.Upgrader

//go:embed index.html
var indexHtml string

var index *template.Template

func init() {
	index = template.Must(template.New("").Parse(indexHtml))
}
