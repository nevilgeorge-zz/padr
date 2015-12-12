// connection.go
package main

import (
  "fmt"
  "github.com/gorilla/websocket"
  "net/http"
)

type connection struct {
  // WebSocket connection
  ws *websocket.Conn

  // Bufferred channel of outbound messages
  send chan []byte

  // Hub connected to
  h *hub
}

// reader method of connection type
func (c *connection) reader() {
  for {
    _, message, err := c.ws.ReadMessage()
    if err != nil {
      fmt.Println("Reading from connection failed.")
      break
    }
    c.h.broadcast <- message
  }
  c.ws.Close()
}

// writer method of connection type
func (c *connection) writer() {
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
  h *hub
}

// method of wsHandler type
func (wsh wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  ws, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    return
  }

  fmt.Println("Adding new connection")
  c := &connection{
    send: make(chan []byte, 256),
    ws: ws,
    h: wsh.h,
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
