// server.go

package main

import (
  "fmt"
  "net"
  "os"
)

type Client struct {
  id int
  conn net.Conn
}

func main() {
  var port string = os.Args[1]
  count := 0

  // listen for incoming connections
  listen, err := net.Listen("tcp", "localhost:" + port)
  if err != nil {
    fmt.Println("Error in listening: ", err.Error())
    os.Exit(1)
  }

  defer listen.Close()
  fmt.Println("Listening on localhost:" + port)

  for {
    conn, err := listen.Accept()
    if err != nil {
      fmt.Println("Error occurred in accepting a connection ", err.Error())
      os.Exit(1)
    }

    newClient := new(Client)
    count += 1
    newClient.id = count
    newClient.conn = conn

  }
}
