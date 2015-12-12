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
  indexTempl = template.Must(template.ParseFiles(filepath.Join(assets, "index.html")))
  h := newHub()
  go h.run()

  http.HandleFunc("/", homeHandler)
  http.Handle("/ws", wsHandler{h: h})
  err := http.ListenAndServe(port, nil)
  if err != nil {
    log.Fatal("ListenAndServe error:", err)
  }
  fmt.Println("Listening for connections on port ", port)
}
