// hub.go
package main

import (
  "encoding/json"
  "fmt"
)

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

  // keep track of cursors
  cursors map[*Connection]map[string]float64
}

type SocketUpdate struct {
  Text string
  Cursors [][]float64
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
    cursors: make(map[*Connection]map[string]float64),
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
      h.cursors[c] = map[string]float64{"start": 0, "end": 0}
      c.send <- h.getSocketUpdate()

    case c := <-h.unregister:
      if _, ok := h.connections[c]; ok {
        h.deleteId(c.id)
        delete(h.connections, c)
        delete(h.cursors, c)
        close(c.send)
      }

    case op := <-h.broadcast:
      h.mergeOperation(op)
      h.cursors[op.sender] = op.selectionRange

      for c := range h.connections {
        if c.id != op.sender.id {
          c.send <- h.getSocketUpdate()
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
  state := h.state
  start := int(op.start)
  addition := []byte(op.chars)
  count := int(op.count)

  switch op.opType {
  case "insert":
    if start == len(state) {
      state = append(state, addition...)
    } else {
      state = append(state[0:start], append(addition, state[start:]...)...)
    }
  case "delete":
    if start == len(state) {
      state = state[0:start - 1]
    } else {
      state = append(state[0: start], state[start + count:]...)
    }
  }

  h.state = state
}

func (h *Hub) getSocketUpdate() []byte {
      ranges := make([][]float64, len(h.cursors))
      for _, val := range h.cursors {
        ranges = append(ranges, []float64{val["start"], val["end"]})
      }
      update := &SocketUpdate{
        Text: string(h.state),
        Cursors: ranges,
      }
      b := convertToBytes(update)
      return b
}

func convertToBytes(update *SocketUpdate) []byte {
  converted, err := json.Marshal(update)
  if err != nil {
    fmt.Println(err)
    return nil
  }
  return converted
}
