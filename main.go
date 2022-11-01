package main

import (
	"embed"
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
	http.Handle("/static/", http.StripPrefix("/static/", static))
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

	go func() {
		var i int
		for {
			<-time.After(time.Second)
			msg := fmt.Sprintf("ping %v", i)
			err = conn.WriteMessage(
				websocket.TextMessage, []byte(msg),
			)
			if err != nil {
				log.Print("pinger closed")
				return
			}
			i++
		}
	}()
	// The event loop
	for {
		_, message, err := conn.ReadMessage()
		if err, ok := err.(*websocket.CloseError); ok {
			log.Print(err)
			break
		}
		log.Printf("Received: %s", message)
	}
	log.Print("socketHandler done")
}

func home(w http.ResponseWriter, r *http.Request) {
	index.ExecuteTemplate(w, "", nil)
}

var upgrader websocket.Upgrader

var (
	//go:embed index.html
	indexHtml string
	index     = template.Must(template.New("").Parse(indexHtml))

	//go:embed ws.js
	assets embed.FS
	static = http.FileServer(http.FS(assets))
)
