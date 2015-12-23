// connection.go
package main

import (
  "encoding/json"
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

type Operation struct {
  // start index of operation
  start float64

  // count of characters affected
  count float64

  // characters added/ [] for deletion
  chars string

  // type of operation (ie. input/delete)
  opType string

  // selection range with start and end keys
  selectionRange map[string]float64

  // Connection that is sending the message
  sender *Connection
}

// reader method of connection type
func (c *Connection) reader() {
  for {
    _, msg, err := c.ws.ReadMessage()

    if err != nil {
      fmt.Println(err)
      fmt.Println("Reading from connection failed.")
      break
    }

    var operation map[string]interface{}
    if err := json.Unmarshal(msg, &operation); err != nil {
      fmt.Println(err)
      continue
    }

    selectionRange := operation["range"].(map[string]interface{})

    // create new operation using type assertion
    op := &Operation{
      start: operation["start"].(float64),
      count: operation["count"].(float64),
      chars: operation["chars"].(string),
      opType: operation["type"].(string),
      selectionRange: map[string]float64{"start": selectionRange["start"].(float64), "end": selectionRange["end"].(float64)},
      sender: c,
    }

    c.h.broadcast <- op

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
