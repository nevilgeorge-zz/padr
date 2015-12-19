// main.go
package main

import (
  // "encoding/json"
  "fmt"
  "log"
  "github.com/gorilla/mux"
  "math/rand"
  "net/http"
  "path/filepath"
  "html/template"
)

var port string = ":8080"
var assets string = "./client"
var indexTempl *template.Template
var sessionTempl *template.Template

// struct to be passed into templates
type TemplateData struct {
  ShortCode string
  Domain string
}

// handle templating for /
func homeHandler(res http.ResponseWriter, req *http.Request) {
  templateVariables := TemplateData{"", req.Host}
  indexTempl.Execute(res, templateVariables)
}

// handle templating for /<ShortCode>
func sessionHandler(res http.ResponseWriter, req *http.Request, shortCode string) {
  templateVariables := TemplateData{shortCode, req.Host}
  sessionTempl.Execute(res, templateVariables)
}


func main() {
  indexTempl = template.Must(template.ParseFiles(filepath.Join(assets, "index.html")))
  sessionTempl = template.Must(template.ParseFiles(filepath.Join(assets, "ws.html")))

  server := NewServer()

  // serve static files at /static/
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./client"))))

  // routes
  homeRoute := Route{
    name: "Home Route",
    method: "GET",
    endpoint: "/",
    handlerFunc: homeHandler,
  }

  createSessionRoute := Route{
    name: "Create Session Route",
    method: "POST",
    endpoint: "/session",
    handlerFunc: func(res http.ResponseWriter, req *http.Request) {
      // use a closure to access the server
      newCode := generateShortCode(6)
      fmt.Println(newCode)

      // create a new hub for this session and run it in a separate goroutine
      hub := newHub()
      hub.shortCode = newCode
      go hub.run()

      ws := wsHandler{h: hub}
      server.AddWsHub(newCode, &ws)

      // create a new route for this new session
      newRoute := Route{
        name: newCode + " Route",
        method: "GET",
        endpoint: "/" + newCode,
        handlerFunc: func(res http.ResponseWriter, req *http.Request) {
          sessionHandler(res, req, newCode)
        },
      }
      server.AddRoute(newRoute)

      // handle socket connections at this new endpoint
      http.Handle("/" + newCode + "/ws", ws)
    },
  }

  getSessionRoute := Route{
    name: "Get Session Route",
    method: "GET",
    endpoint: "/{shortCode}",
    handlerFunc: func(res http.ResponseWriter, req *http.Request) {
      // grab params from request
      vars := mux.Vars(req)
      shortCode := vars["shortCode"]
      ws := server.wsHubs[shortCode]
      if ws == nil {
        http.Redirect(res, req, "/", http.StatusFound)
      } else {
        sessionHandler(res, req, shortCode)
      }
    },
  }

  // add newly created routes
  server.AddRoute(homeRoute)
  server.AddRoute(createSessionRoute)
  server.AddRoute(getSessionRoute)

  http.Handle("/", server.router)

  err := http.ListenAndServe(port, nil)
  if err != nil {
    log.Fatal("ListenAndServe error:", err)
  }
}

func generateShortCode(length int) string {
  const possibleBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

  code := make([]byte, length)
  for i := 0; i < length; i++ {
    code[i] = possibleBytes[rand.Intn(len(possibleBytes))]
  }

  return string(code)
}
