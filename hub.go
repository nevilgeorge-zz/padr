// hub.go
package main

type hub struct {
  // Registered connections
  connections map[*connection]bool

  // Inbound messages from connections
  broadcast chan []byte

  // Register requests from the connections
  register chan *connection

  // Unregister requests from the connections
  unregister chan *connection
}

// constructor for hub struct
func newHub() *hub {
  hub := hub{
    connections: make(map[*connection]bool),
    broadcast: make(chan []byte),
    register: make(chan *connection),
    unregister: make(chan *connection),
  }

  return &hub
}

// run method for hub type
func (h *hub) run() {
  for {
    select {
    case c := <-h.register:
      h.connections[c] = true

    case c := <-h.unregister:
      if _, ok := h.connections[c]; ok {
        delete(h.connections, c)
        close(c.send)
      }

    case m := <-h.broadcast:
      for c := range h.connections {
        select {
        case c.send <- m:
        default:
          delete(h.connections, c)
          close(c.send)
        }
      }
    }
  }
}
