// main.go
package main

import (
  "fmt"
  "log"
  "net/http"
  "path/filepath"
  "text/template"
)

var port string = ":8080"
var assets string = "./client"
var indexTempl *template.Template

func homeHandler(res http.ResponseWriter, req *http.Request) {
  indexTempl.Execute(res, req.Host)
}

func main() {
  indexTempl = template.Must(template.ParseFiles(filepath.Join(assets, "ws.html")))
  h := newHub()
  go h.run()

  // serve static files at /static/
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./client"))))

  http.HandleFunc("/", homeHandler)
  http.Handle("/ws", wsHandler{h: h})
  err := http.ListenAndServe(port, nil)
  if err != nil {
    log.Fatal("ListenAndServe error:", err)
  }
  fmt.Println("Listening for connections on port ", port)
}
