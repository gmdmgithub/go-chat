package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// client represents a single chatting user.
type client struct {
	// socket is the web socket for this client.
	socket *websocket.Conn
	// send is a channel on which messages are sent.
	send chan *message
	// room is the room this client is chatting in.
	room *room
	// userData holds information about the user- if it shuld be a googleUser??
	userData map[string]interface{}
}

func (c *client) read() {
	for {
		var msg *message
		// if _, msg, err := c.socket.ReadMessage(); err == nil {
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			// if picture, ok := c.userData["picture"]; ok {
			// 	msg.Picture = picture.(string)
			// }
			msg.Picture, _ = c.room.avatar.GetAvatarURL(c)
			c.room.forward <- msg
			log.Printf("client read %v", c.userData)
			log.Printf("message read %v", msg)
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		// if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
		log.Printf("message write %v", msg)
	}
	c.socket.Close()
}
