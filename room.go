package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gmdmgithub/chat/trace"
	"github.com/gorilla/websocket"
)

type room struct {
	// forward is a channel that holds incoming messages
	// that should be forwarded to the other clients.
	forward chan *message
	// join is a channel for clients wishing to join the room.
	join chan *client
	// leave is a channel for clients wishing to leave the room.
	leave chan *client
	// clients holds all current clients in this room.
	clients map[*client]bool
	// tracer will receive trace information of activity
	// in the room.
	tracer trace.Tracer
}

// run - infinitive method to handle three defined channels
func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// joining
			r.clients[client] = true
			r.tracer.Trace("New client joined")
		case client := <-r.leave:
			// leaving
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client left")
		case msg := <-r.forward:
			//r.tracer.Trace("Message received: ", msg)
			r.tracer.Trace("Message received: ")
			// forward message to all clients
			for client := range r.clients {
				select {
				case client.send <- msg:
					// send the message
					r.tracer.Trace(" ... sent to client")
				default:
					// failed to send
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace("failed to send... cleaned up client")
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

// ServeHTPP for the room struct allows to triger ServeHTTP - for case room is send as an interface of http.Handler
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println("room request started")
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Printf("ServeHTTP FATAL PROBLEM: %v\n", err)
		http.Redirect(w, req, "/chat", http.StatusTemporaryRedirect)
		return
	}
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie FATAL PROBLEM: ", err)
		return
	}

	usrData, err := base64.StdEncoding.DecodeString(authCookie.Value)
	if err != nil {
		log.Fatal("Failed to get user data from cookie FATAL PROBLEM: ", err)
		return
	}

	log.Println("room serve")

	// var user *googleUser
	// err = json.Unmarshal(usrData, &user)
	var rawData map[string]interface{}
	err = json.Unmarshal(usrData, &rawData)

	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: rawData,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}

// newRoom makes a new room that is ready to go.
func newRoom() *room {
	return &room{
		// forward: make(chan []byte),
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		// tracer:  trace.New(os.Stdout),
		tracer: trace.Off(),
	}
}
