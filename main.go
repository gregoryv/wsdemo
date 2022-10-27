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
	http.HandleFunc("/ws.js", wsjs)

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

	go func() {
		var i int
		for {
			<-time.After(2 * time.Second)
			msg := fmt.Sprintf("ping %v", i)
			err = conn.WriteMessage(
				websocket.TextMessage, []byte(msg),
			)
			if err != nil {
				log.Println("Error during message writing:", err)
			}
			i++
		}
	}()
	// The event loop
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			break
		}
		log.Printf("Received: %s", message)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	index.ExecuteTemplate(w, "", nil)
}

func wsjs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/javascript")
	w.Write(wsJS)
}

var upgrader websocket.Upgrader

var (
	//go:embed index.html
	indexHtml string
	index     = template.Must(template.New("").Parse(indexHtml))

	//go:embed ws.js
	wsJS []byte
)
