// main.go
package main

import (
  "go/build"
  "log"
  "net/http"
  "path/filepath"
  "text/template"
)

var port string = ":8080"
var assets string = "./client"
var homeTempl *template.Template

func homeHandler(res http.ResponseWriter, req *http.Request) {
  homeTempl.Execute(res, req.Host)
}

func main() {
  homeTempl = template.Must(template.ParseFiles(filepath.Join(*assets, "home.html")))
  h := newHub()
  go h.run()

  http.HandleFunc("/", homeHandler)
  http.Handle("/ws", wsHandler{h: h})
  if err := http.ListenAndServe(*port, nil); err != nil {
    log.Fatal("ListenAndServe error:", err)
  }
}
