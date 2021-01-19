package main

import (
	"fmt"
	"log"
	"net/http"
	"io"
	"io/ioutil"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func ws(w http.ResponseWriter, req *http.Request){
	c, err := upgrader.Upgrade(w, req, nil) // Upgrade TCP request to websocket
	if err != nil{
		log.Print("Upgrade:", err)
		return
	}
	defer c.Close()
	for{
		mt, message, err := c.ReadMessage() // Read from websocket
		if err != nil{
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)
		if err != nil{
			log.Println("write:", err)
			return
		}
		err = c.WriteMessage(mt, message) // Write to websocket
		if err != nil{
			log.Println("write:", err)
			return
		}
	}
}

func main()  {
	clients := make(map[string]bool)
	
	// Handler for joining the chat
	http.HandleFunc("/join", func(w http.ResponseWriter, req *http.Request){
		req.ParseForm()
		fmt.Println(req.Host)
		clients[req.PostFormValue("name")] = true
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		io.WriteString(w, "Registered successfully")
	})

	// Render index page
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request){
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		content, _ := ioutil.ReadFile("./templates/index.html")
		w.Write(content)
	})

	// Respond to chat messages
	http.HandleFunc("/echo", ws)

	// Start server
	log.Fatal(http.ListenAndServe(":8090", nil))
}