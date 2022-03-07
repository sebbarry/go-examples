package main

import (
    "fmt"
    "net"
    "os"
    "strings"
)

const (
    CONN_HOST = "localhost"
    CONN_PORT = "3333"
    CONN_TYPE = "tcp"
)



func main() {
    startRequest()
}

func startRequest() {
    // Listen for incoming connections.
    l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()
    fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        // Handle connections in a new goroutine.
        go handleRequest(conn)
    }
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
    //defered function to run when the loop is finished
  defer func() {
      if err := conn.Close(); err != nil {
          fmt.Println("error closing connection: ", err)
      }
  }()

  if _,err := conn.Write([]byte ("Connected...\nUsage: GET <value>")); err != nil {
      fmt.Println("error writing:", err)
      return
  }

  //this is the conection loop to keep the connection with the client.
  for {
      //make a buffer slice of 1024 bytes.
      buf := make([]byte, 1024)
      // Read the incoming connection into the buffer.
      n, err := conn.Read(buf)
      if n == 0 || err != nil {
          fmt.Println("Connection read error: ", err)
          return
      }
      fmt.Println(string(buf[0:n])) //convert to string
      cmd, _ := parseCommand(string(buf[0:n]))
      if cmd == "" {
          if _,err := conn.Write([]byte("Invalid command\n")); err !=nil {
              fmt.Println("failed to write:", err)
              return
          }
          continue
      }
      switch strings.ToUpper(cmd){
      case "GET":
          if _,err := conn.Write([]byte("you sent get.")); err != nil {
              fmt.Println("failed to write to client", err)
          }
          continue
      }
  }
}

func parseCommand(cmdLine string) (cmd, param string) {
    parts := strings.Split(cmdLine, " ")
    if len(parts) != 2 {
        return "", ""
    }
    cmd = strings.TrimSpace(parts[0])
    param = strings.TrimSpace(parts[1])
    return
}
