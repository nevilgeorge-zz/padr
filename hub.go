// hub.go
package main

type hub struct {
  // Registered connections
  connections map[*connection]bool

  // Inbound messages from connections
  broadcast chan *message

  // Register requests from the connections
  register chan *connection

  // Unregister requests from the connections
  unregister chan *connection

  // list of assigned ids
  connection_ids []int

  // map that holds state for each session
  state map[int][]byte
}

// constructor for hub struct
func newHub() *hub {
  hub := hub{
    connections: make(map[*connection]bool),
    broadcast: make(chan *message),
    register: make(chan *connection),
    unregister: make(chan *connection),
    connection_ids: make([]int, 0),
    state: make(map[int][]byte),
  }

  return &hub
}

// run method for hub type
func (h *hub) run() {
  for {
    select {
    case c := <-h.register:
      h.connections[c] = true
      c.id = h.get_next_id()
      c.send <- h.state[c.session_id]

    case c := <-h.unregister:
      if _, ok := h.connections[c]; ok {
        h.delete_id(c.id)
        delete(h.connections, c)
        close(c.send)
      }

    case m := <-h.broadcast:
      h.update_state(m.sender.session_id, m.text)
      for c := range h.connections {
        if c.id != m.sender.id {
          c.send <- m.text
        }
      }
    }
  }
}

// function to allocate new connection id
func (h *hub) get_next_id() int {

  ids := make(map[int]bool)

  for i := 0; i < len(h.connection_ids); i++ {
    ids[h.connection_ids[i]] = true
  }

  for i := 0; i < len(h.connection_ids); i++ {
    if ids[i + 1] == false {
      return i + 1
    }
  }

  h.connection_ids = append(h.connection_ids, len(h.connection_ids) + 1)
  return h.connection_ids[len(h.connection_ids) - 1]
}

// function to deallocate connection_id
func (h *hub) delete_id(id int) {

  index := -1
  // find index of id
  for i, value := range h.connection_ids {
    if value == id {
      index = i
    }
  }

  if index == -1 {
    return
  }

  h.connection_ids = append(h.connection_ids[:index], h.connection_ids[index + 1:]...)
}

func (h *hub) update_state(session_id int, new_state []byte) {
  h.state[session_id] = new_state
}
