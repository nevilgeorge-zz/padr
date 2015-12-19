// connection.go
package main

import (
  "fmt"
  "github.com/gorilla/websocket"
  "net/http"
)

type Connection struct {
  // WebSocket connection
  ws *websocket.Conn

  // Bufferred channel of outbound messages
  send chan []byte

  // Hub connected to
  h *Hub

  // Assigned id
  id int

  // session id
  sessionId int
}

type Message struct {
  // Text of string being sent
  text []byte

  // Connection that is sending the message
  sender *Connection
}

// reader method of connection type
func (c *Connection) reader() {
  for {
    _, msg, err := c.ws.ReadMessage()

    newMessage := &Message{text: msg, sender: c}

    if err != nil {
      fmt.Println(err)
      fmt.Println("Reading from connection failed.")
      break
    }
    c.h.broadcast <- newMessage
  }
  c.ws.Close()
}

// writer method of Connection type
func (c *Connection) writer() {
  for message := range c.send {
    err := c.ws.WriteMessage(websocket.TextMessage, message)
    if err != nil {
      fmt.Println("Writing to connection failed.")
      break
    }
  }
  c.ws.Close()
}

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

type wsHandler struct {
  h *Hub
}

// method of wsHandler type
func (wsh wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  ws, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    return
  }

  fmt.Println("Adding new connection")
  c := &Connection{
    send: make(chan []byte, 256),
    ws: ws,
    h: wsh.h,
    id: -1,
    sessionId: 1,
  }
  c.h.register <- c

  defer func() {
    c.h.unregister <- c
  }()

  // new goroutine for writer
  go c.writer()
  // use this goroutine for reader
  c.reader()
}
