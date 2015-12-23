// hub.go
package main

type Hub struct {
  // Registered connections
  connections map[*Connection]bool

  // Inbound messages from connections
  broadcast chan *Operation

  // Register requests from the connections
  register chan *Connection

  // Unregister requests from the connections
  unregister chan *Connection

  // list of assigned ids
  connectionIds []int

  // map that holds state for current session
  state []byte

  // shortcode associated with this hub
  shortCode string
}

// constructor for hub struct
func newHub() *Hub {
  hub := Hub{
    connections: make(map[*Connection]bool),
    broadcast: make(chan *Operation),
    register: make(chan *Connection),
    unregister: make(chan *Connection),
    connectionIds: make([]int, 0),
    state: make([]byte, 0),
    shortCode: "",
  }

  return &hub
}

// run method for hub type
func (h *Hub) run() {
  for {
    select {
    case c := <-h.register:
      c.id = h.getNextId()
      h.connections[c] = true
      c.send <- h.state

    case c := <-h.unregister:
      if _, ok := h.connections[c]; ok {
        h.deleteId(c.id)
        delete(h.connections, c)
        close(c.send)
      }

    case op := <-h.broadcast:
      h.mergeOperation(op)
      for c := range h.connections {
        if c.id != op.sender.id {
          c.send <- h.state
        }
      }
    }
  }
}

// function to allocate new connection id
func (h *Hub) getNextId() int {

  ids := make(map[int]bool)

  for i := 0; i < len(h.connectionIds); i++ {
    ids[h.connectionIds[i]] = true
  }

  for i := 0; i < len(h.connectionIds); i++ {
    if ids[i + 1] == false {
      return i + 1
    }
  }

  h.connectionIds = append(h.connectionIds, len(h.connectionIds) + 1)
  return h.connectionIds[len(h.connectionIds) - 1]
}

// function to deallocate connection_id
func (h *Hub) deleteId(id int) {

  index := -1
  // find index of id
  for i, value := range h.connectionIds {
    if value == id {
      index = i
    }
  }

  if index == -1 {
    return
  }

  h.connectionIds = append(h.connectionIds[:index], h.connectionIds[index + 1:]...)
}

func (h *Hub) mergeOperation(op *Operation) {
  h.state = []byte("Hello")
}
