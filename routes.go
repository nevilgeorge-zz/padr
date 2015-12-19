// routes.go
package main

import (
  "github.com/gorilla/mux"
  "net/http"
)

type Route struct {

  // name the route
  name string

  // 'GET' vs 'POST', etc
  method string

  // "/hello"
  endpoint string

  // function to handle the request
  handlerFunc http.HandlerFunc
}

type Server struct {
  // gorilla request multiplexer
  router *mux.Router

  // one hub per session
  wsHubs map[string]*wsHandler

}

func NewServer() *Server {

  router := mux.NewRouter()
  server := Server{
    router: router,
    wsHubs: make(map[string]*wsHandler),
  }

  return &server
}

// add a hub to the router
func (server *Server) AddWsHub(shortCode string, ws *wsHandler) {
  server.wsHubs[shortCode] = ws
}

// add a single route to router
func (server *Server) AddRoute(route Route) {
  server.router.
    Methods(route.method).
    Path(route.endpoint).
    Name(route.name).
    Handler(route.handlerFunc)

}
