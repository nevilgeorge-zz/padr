// hub.go
package main

type Hub struct {
  // Registered connections
  connections map[*Connection]bool

  // Inbound messages from connections
  broadcast chan *Message

  // Register requests from the connections
  register chan *Connection

  // Unregister requests from the connections
  unregister chan *Connection

  // list of assigned ids
  connectionIds []int

  // map that holds state for each session
  state map[int][]byte

  // shortcode associated with this hub
  shortCode string
}

// constructor for hub struct
func newHub() *Hub {
  hub := Hub{
    connections: make(map[*Connection]bool),
    broadcast: make(chan *Message),
    register: make(chan *Connection),
    unregister: make(chan *Connection),
    connectionIds: make([]int, 0),
    state: make(map[int][]byte),
    shortCode: "",
  }

  return &hub
}

// run method for hub type
func (h *Hub) run() {
  for {
    select {
    case c := <-h.register:
      h.connections[c] = true
      c.id = h.getNextId()
      c.send <- h.state[c.sessionId]

    case c := <-h.unregister:
      if _, ok := h.connections[c]; ok {
        h.deleteId(c.id)
        delete(h.connections, c)
        close(c.send)
      }

    case m := <-h.broadcast:
      h.updateState(m.sender.sessionId, m.text)
      for c := range h.connections {
        if c.id != m.sender.id {
          c.send <- m.text
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

func (h *Hub) updateState(sessionId int, newState[]byte) {
  h.state[sessionId] = newState
}
